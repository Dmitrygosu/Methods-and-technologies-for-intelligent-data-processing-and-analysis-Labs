import sys
import os
import multiprocessing
import torch
# Add the project root to sys.path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from sb3_contrib import MaskablePPO
from stable_baselines3.common.vec_env import SubprocVecEnv
from stable_baselines3.common.utils import set_random_seed
from stable_baselines3.common.callbacks import CheckpointCallback, BaseCallback
from core.env import FrisianCheckersEnv

class CustomMetricsCallback(BaseCallback):
    def __init__(self, verbose=0):
        super(CustomMetricsCallback, self).__init__(verbose)
        self.win_count = 0
        self.episode_count = 0

    def _on_step(self) -> bool:
        # Check for info and dones from all environments
        infos = self.locals.get("infos")
        dones = self.locals.get("dones")
        if infos is not None and dones is not None:
            for info, done in zip(infos, dones):
                if done:
                    self.episode_count += 1
                    if info.get("is_win", 0):
                        self.win_count += 1
                    
                    # Log pieces stats (at the end of episode)
                    self.logger.record("env/pieces_lost", info.get("pieces_lost", 0))
                    self.logger.record("env/pieces_captured", info.get("pieces_captured", 0))
                    
        # Log winrate periodically
        if self.episode_count > 0:
            self.logger.record("env/win_rate", self.win_count / self.episode_count)
            
        return True

def make_env(rank, seed=0):
    def _init():
        env = FrisianCheckersEnv()
        env.reset(seed=seed + rank)
        return env
    set_random_seed(seed)
    return _init

def train():
    num_cpu = 12 
    device = "cuda" if torch.cuda.is_available() else "cpu"
    print(f"Training (Deep MLP): Threads: {num_cpu} | Device: {device}")
    
    env = SubprocVecEnv([make_env(i) for i in range(num_cpu)])
    
    # Deep MLP Policy: [512, 512, 256] architecture
    policy_kwargs = dict(
        net_arch=dict(pi=[512, 512, 256], vf=[512, 512, 256])
    )
    
    model = MaskablePPO(
        "MlpPolicy", 
        env, 
        policy_kwargs=policy_kwargs,
        verbose=1, 
        learning_rate=5e-4, 
        n_steps=2048,
        batch_size=512,
        n_epochs=10,
        gamma=0.999,
        ent_coef=0.01,
        clip_range=0.2,
        max_grad_norm=0.5,
        device=device,
        tensorboard_log=os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "ppo_checkers_tensorboard")
    )
    
    checkpoint_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "ai", "checkpoints")
    os.makedirs(checkpoint_dir, exist_ok=True)

    checkpoint_callback = CheckpointCallback(
        save_freq=250000 // num_cpu,
        save_path=checkpoint_dir,
        name_prefix="poddavki_rl_model",
        save_replay_buffer=True,
        save_vecnormalize=True,
    )

    # Save initial model so parallel envs can start Self-Play immediately
    model.save("ai/poddavki_rl_model")
    print("Initial model saved for Self-Play startup.")

    custom_callback = CustomMetricsCallback()
    
    # Combined callbacks
    from stable_baselines3.common.callbacks import CallbackList
    callbacks = CallbackList([checkpoint_callback, custom_callback])

    print("Starting Training... 3,100,000 steps.")
    model.learn(total_timesteps=3100000, callback=callbacks)
    
    model.save("ai/poddavki_rl_model")
    print("Model saved to ai/poddavki_rl_model.zip")

if __name__ == "__main__":
    multiprocessing.freeze_support()
    train()
