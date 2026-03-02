# Makefile for Frisian Checkers

PYTHON = python
MAIN = gui/main_window.py
TRAIN = scripts/train.py
BENCHMARK = scripts/benchmark.py
TENSORBOARD_DIR = ppo_checkers_tensorboard

.PHONY: run train benchmark tensorboard clean help

help:
	@echo "Available commands:"
	@echo "  make run         - Start the game GUI"
	@echo "  make train       - Start RL agent training"
	@echo "  make benchmark   - Run AI performance benchmarking"
	@echo "  make tensorboard - Start TensorBoard to track training progress"
	@echo "  make clean       - Remove temporary files and __pycache__"

run:
	$(PYTHON) $(MAIN)

train:
	$(PYTHON) $(TRAIN)

benchmark:
	$(PYTHON) $(BENCHMARK)

tensorboard:
	tensorboard --logdir=$(TENSORBOARD_DIR)

clean:
	rm -rf **/__pycache__
	rm -rf .pytest_cache
	rm -rf *.spec
	rm -rf build/
	rm -rf dist/
