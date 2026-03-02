import os
import numpy as np
from core.engine import BOARD_SIZE, WHITE_MAN, BLACK_MAN, WHITE_KING, BLACK_KING

try:
    from sb3_contrib import MaskablePPO
    RL_AVAILABLE = True
except (ImportError, OSError) as e:
    print(f"Warning: RL libraries (sb3-contrib/torch) could not be loaded: {e}")
    RL_AVAILABLE = False

class RLAgent:
    def __init__(self, model_path="ai/poddavki_rl_model.zip"):
        self.model = None
        if not RL_AVAILABLE: return

        # Try to find the best model
        base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        
        candidates = [
            os.path.join(base_dir, model_path),
            os.path.join(base_dir, "ai/poddavki_rl_model.zip"),
            # Check for latest models
        ]
        
        # Add checkpoints if they exist (PRIORITIZE CHECKPOINTS)
        checkpoint_dir = os.path.join(base_dir, "ai/checkpoints")
        if os.path.exists(checkpoint_dir):
            ckpt_files = [f for f in os.listdir(checkpoint_dir) if f.endswith(".zip")]
            if ckpt_files:
                # Sort by filename (which includes timestep)
                ckpt_files.sort(key=lambda x: int(''.join(filter(str.isdigit, x)) or 0), reverse=True)
                # INSERT AT BEGINNING to prioritize the latest evolved model
                candidates.insert(0, os.path.join(checkpoint_dir, ckpt_files[0]))

        for path in candidates:
            if os.path.exists(path):
                try:
                    self.model = MaskablePPO.load(path)
                    print(f"Success: Loaded RL model from {path}")
                    return
                except Exception as e:
                    print(f"Error loading {path}: {e}")

        print("WARNING: No RL model found. Falling back to Random.")


    def _pos_to_idx(self, r, c):
        return (r * 10 + c) // 2

    def _idx_to_pos(self, idx):
        r = (idx * 2) // 10
        c = (idx * 2) % 10
        if (r + c) % 2 == 0: c += 1
        return (r, c)

    def get_move(self, engine, player):
        legal_moves = engine.get_legal_moves(player)
        if not legal_moves: return None
        if not self.model:
            import random
            return random.choice(legal_moves)
        
        obs = engine.get_state()
        if player == BLACK_MAN:
            obs = np.flip(obs)
            new_obs = np.zeros_like(obs)
            new_obs[obs == WHITE_MAN] = BLACK_MAN
            new_obs[obs == BLACK_MAN] = WHITE_MAN
            new_obs[obs == WHITE_KING] = BLACK_KING
            new_obs[obs == BLACK_KING] = WHITE_KING
            obs = new_obs
        
        # Prepare features for model input
        agent_pieces = sum(1 for p in obs.flatten() if p in [WHITE_MAN, WHITE_KING])
        opp_pieces = sum(1 for p in obs.flatten() if p in [BLACK_MAN, BLACK_KING])
        norm_agent = agent_pieces / 20.0
        norm_opp = opp_pieces / 20.0
        legal_for_agent = engine.get_legal_moves(player)
        is_mandatory = 1.0 if any(len(m[2]) > 0 for m in legal_for_agent) else 0.0
        
        # Turn phase heuristic
        turn_phase = 0.5 
        
        # Flattened for MlpPolicy: (100,) + (4,) extra
        flat_board = (obs.astype(np.float32) / 4.0).flatten()
        obs_normalized = np.concatenate([flat_board, [norm_agent, norm_opp, is_mandatory, turn_phase]])
        
        mask = np.zeros(2500, dtype=bool)
        for move in legal_moves:
            s_idx = self._pos_to_idx(*move[0])
            e_idx = self._pos_to_idx(*move[1])
            if player == BLACK_MAN:
                s_idx, e_idx = 49 - s_idx, 49 - e_idx
            idx = s_idx * 50 + e_idx
            if 0 <= idx < 2500: mask[idx] = True
        
        if not np.any(mask):
            import random
            return random.choice(legal_moves)
        
        action, _ = self.model.predict(obs_normalized, action_masks=mask, deterministic=True)
        
        c_start, c_end = action // 50, action % 50
        raw_s, raw_e = (c_start, c_end) if player == WHITE_MAN else (49 - c_start, 49 - c_end)
        sp = self._idx_to_pos(raw_s)
        ep = self._idx_to_pos(raw_e)
        
        for move in legal_moves:
            if move[0] == sp and move[1] == ep:
                return move
        
        print(f"Agent fallback: Model action validation failed.")
        import random
        return random.choice(legal_moves)
