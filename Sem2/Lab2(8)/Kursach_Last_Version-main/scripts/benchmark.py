import random
import time
import sys
import os

# Add the project root to sys.path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from core.engine import FrisianCheckersLogic, WHITE_MAN, BLACK_MAN

class RandomBot:
    def get_move(self, engine, player):
        moves = engine.get_legal_moves(player)
        return random.choice(moves) if moves else None

class Tournament:
    def __init__(self, games=100):
        self.games = games
        self.results = {"White": 0, "Black": 0, "Draw": 0}

    def run(self, bot1, bot2):
        print(f"Starting tournament: {self.games} games...")
        for i in range(self.games):
            engine = FrisianCheckersLogic()
            moves_count = 0
            while not engine.winner and moves_count < 500:
                player = engine.turn
                bot = bot1 if player == WHITE_MAN else bot2
                move = bot.get_move(engine, player)
                if not move:
                    break
                engine.make_move(move)
                moves_count += 1
            
            if engine.winner:
                if "White" in engine.winner: self.results["White"] += 1
                else: self.results["Black"] += 1
            else:
                self.results["Draw"] += 1
            
            if (i + 1) % 10 == 0:
                print(f"Played {i+1}/{self.games} games...")

        print("\n--- Tournament Results ---")
        for k, v in self.results.items():
            print(f"{k}: {v} ({(v/self.games)*100:.1f}%)")

if __name__ == "__main__":
    from ai.agent import RLAgent
    
    # Try to load RL Agent for benchmarking
    bot_rl = RLAgent()
    bot_rand = RandomBot()
    
    print("--- Benchmark: RL as WHITE vs Random ---")
    t1 = Tournament(games=100)
    t1.run(bot_rl, bot_rand)
    
    print("\n--- Benchmark: RL as BLACK vs Random ---")
    t2 = Tournament(games=100)
    t2.run(bot_rand, bot_rl) # RL is bot2 (Black)
