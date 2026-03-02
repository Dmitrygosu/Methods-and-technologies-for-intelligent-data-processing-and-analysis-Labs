import gymnasium as gym
from gymnasium import spaces
import numpy as np
import random
from core.engine import FrisianCheckersLogic, WHITE_MAN, BLACK_MAN, EMPTY, BOARD_SIZE, WHITE_KING, BLACK_KING

class FrisianCheckersEnv(gym.Env):
    metadata = {"render_modes": ["human", "rgb_array"], "render_fps": 4}

    def __init__(self, render_mode=None):
        super(FrisianCheckersEnv, self).__init__()
        self.logic = FrisianCheckersLogic()
        # Observation: Flattened 10x10 board (100) + [agent_pieces, opp_pieces, is_mandatory_capture, turn_phase]
        self.observation_space = spaces.Box(low=0.0, high=1.0, shape=(104,), dtype=np.float32)
        self.action_space = spaces.Discrete(2500)
        self.render_mode = render_mode
        self.agent_color = WHITE_MAN
        self.opponent_bot = None # To be lazy-loaded
        self.total_steps = 0
        self.game_steps = 0  # Per-game step counter for turn_phase

    def _pos_to_idx(self, r, c):
        return (r * 10 + c) // 2

    def _idx_to_pos(self, idx):
        r = (idx * 2) // 10
        c = (idx * 2) % 10
        if (r + c) % 2 == 0: c += 1
        return (r, c)

    def _get_canonical_board(self):
        board = self.logic.board.copy()
        if self.agent_color == BLACK_MAN:
            board = np.flip(board)
            new_board = np.zeros_like(board)
            new_board[board == WHITE_MAN] = BLACK_MAN
            new_board[board == BLACK_MAN] = WHITE_MAN
            new_board[board == WHITE_KING] = BLACK_KING
            new_board[board == BLACK_KING] = WHITE_KING
            board = new_board
            
        # Agent always perceives pieces as WHITE in canonical format
        agent_pieces = sum(1 for p in board.flatten() if p in [WHITE_MAN, WHITE_KING])
        opp_pieces = sum(1 for p in board.flatten() if p in [BLACK_MAN, BLACK_KING])
        
        # Normalize piece counts (max 20)
        norm_agent = agent_pieces / 20.0
        norm_opp = opp_pieces / 20.0
        
        # Check if agent has a mandatory capture right now
        legal_moves = self.logic.get_legal_moves(self.agent_color)
        is_mandatory = 1.0 if any(len(m[2]) > 0 for m in legal_moves) else 0.0
        
        # Turn phase (how deep into the CURRENT GAME we are)
        turn_phase = min(1.0, self.game_steps / 100.0)
        
        flat_board = (board.astype(np.float32) / 4.0).flatten()
        return np.concatenate([flat_board, [norm_agent, norm_opp, is_mandatory, turn_phase]])

    def action_masks(self):
        mask = np.zeros(2500, dtype=bool)
        if self.logic.winner:
            mask[0] = True
            return mask

        # Masks are always generated for the AGENT's color
        legal_moves = self.logic.get_legal_moves(self.agent_color)
        for move in legal_moves:
            s_idx = self._pos_to_idx(*move[0])
            e_idx = self._pos_to_idx(*move[1])
            if self.agent_color == BLACK_MAN:
                s_idx, e_idx = 49 - s_idx, 49 - e_idx
            idx = s_idx * 50 + e_idx
            if 0 <= idx < 2500: mask[idx] = True
        
        if not np.any(mask): mask[0] = True
        return mask

    def _opponent_move(self):
        if self.logic.winner: return
        # Periodically reload opponent bot to ensure dynamic self-play
        self.total_steps += 1
        if self.total_steps % 5000 == 0:
            self.opponent_bot = None

        opponent_color = WHITE_MAN if self.agent_color == BLACK_MAN else BLACK_MAN
        legal_moves = self.logic.get_legal_moves(opponent_color)
        if not legal_moves: return
        
        # 50% Pseudo-Self-Play: Use the current RL model if available
        if random.random() < 0.5:
            try:
                if self.opponent_bot is None:
                    from ai.agent import RLAgent
                    self.opponent_bot = RLAgent()
                
                # Check if model loaded successfully
                if self.opponent_bot.model is not None:
                    move = self.opponent_bot.get_move(self.logic, opponent_color)
                    if move:
                        self.logic.make_move(move)
                        return
            except Exception:
                pass # Fallback to random if model not available or error

        # Fallback: Random Move
        move = random.choice(legal_moves)
        self.logic.make_move(move)

    def _evaluate_state(self):
        white = sum(1 for p in self.logic.board.flatten() if self.logic.is_white(p))
        black = sum(1 for p in self.logic.board.flatten() if self.logic.is_black(p))
        if self.agent_color == WHITE_MAN:
            return white, black
        else:
            return black, white

    def reset(self, seed=None, options=None):
        super().reset(seed=seed)
        self.logic = FrisianCheckersLogic()
        self.game_steps = 0  # Reset per-game counter
        
        # Assign Agent's color randomly
        if seed is not None:
            random.seed(seed)
        self.agent_color = random.choice([WHITE_MAN, BLACK_MAN])
        
        # If agent is Black, Opponent (White) moves first
        if self.agent_color == BLACK_MAN:
            self._opponent_move()
            
        return self._get_canonical_board(), {}

    def step(self, action):
        if self.logic.winner:
            return self._get_canonical_board(), 0, True, False, {}

        self.game_steps += 1
        c_start, c_end = action // 50, action % 50
        raw_s, raw_e = (c_start, c_end) if self.agent_color == WHITE_MAN else (49 - c_start, 49 - c_end)
        s_pos, e_pos = self._idx_to_pos(raw_s), self._idx_to_pos(raw_e)
        
        agent_before, opp_before = self._evaluate_state()
        
        legal_moves = self.logic.get_legal_moves(self.agent_color)
        move = next((m for m in legal_moves if m[0] == s_pos and m[1] == e_pos), None)
        
        reward = 0.0
        info = {"is_win": 0, "pieces_lost": 0, "pieces_captured": 0}
        
        if move:
            # 1. Agent Moves
            self.logic.make_move(move)
            
            # Measure state AFTER agent's move but BEFORE opponent's move
            agent_mid, opp_mid = self._evaluate_state()
            
            # Movement rewards
            # Calculate piece differences
            agent_lost = agent_before - agent_mid
            agent_captured = opp_before - opp_mid
            
            # +20.0 for losing a piece (Sacrifice), -15.0 for capturing pieces
            reward += (agent_lost * 20.0) - (agent_captured * 15.0)
            
            info["pieces_lost"] = agent_lost
            info["pieces_captured"] = agent_captured
            
            # 2. Opponent Moves (if game not over)
            if not self.logic.winner:
                # Check for forced captures for opponent
                opponent_color = WHITE_MAN if self.agent_color == BLACK_MAN else BLACK_MAN
                opp_moves = self.logic.get_legal_moves(opponent_color)
                if opp_moves and any(len(m[2]) > 0 for m in opp_moves):
                    reward += 10.0  # Strategic reward
                    
                self._opponent_move()
            
            if self.logic.winner:
                win_color = "White" if self.agent_color == WHITE_MAN else "Black"
                if win_color in self.logic.winner:
                    reward += 1000.0  # Win reward
                    info["is_win"] = 1
                else:
                    reward -= 1000.0  # Loss penalty
        else:
            # CHECK FOR WIN BY BLOCK (No moves left)
            # In Poddavki, if you have no moves, you WIN.
            if not legal_moves and self.logic.winner:
                win_color = "White" if self.agent_color == WHITE_MAN else "Black"
                if win_color in self.logic.winner:
                    reward = 1000.0
                    info["is_win"] = 1
                    return self._get_canonical_board(), reward, True, False, info

            # Otherwise, it's either an invalid action index or a crash
            reward = -20.0
            return self._get_canonical_board(), reward, True, False, info
            
        return self._get_canonical_board(), reward, bool(self.logic.winner), False, info

    def render(self): pass
