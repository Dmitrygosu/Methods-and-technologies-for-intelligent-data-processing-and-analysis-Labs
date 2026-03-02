import numpy as np
from typing import List, Tuple, Optional

# Constants
EMPTY = 0
WHITE_MAN = 1
BLACK_MAN = 2
WHITE_KING = 3
BLACK_KING = 4

BOARD_SIZE = 10

class FrisianCheckersLogic:
    """
    Core engine for Frisian Checkers (Giveaway/Poddavki variant).
    Implements 10x10 board logic with diagonal and orthogonal captures.
    """
    def __init__(self):
        self.board = self._create_initial_board()
        self.turn = WHITE_MAN  # White starts (1 or 3)
        self.winner = None

    def _create_initial_board(self) -> np.ndarray:
        board = np.zeros((BOARD_SIZE, BOARD_SIZE), dtype=int)
        # Black pieces (rows 0-3)
        for r in range(4):
            for c in range(BOARD_SIZE):
                if (r + c) % 2 != 0:
                    board[r, c] = BLACK_MAN
        # White pieces (rows 6-9)
        for r in range(6, BOARD_SIZE):
            for c in range(BOARD_SIZE):
                if (r + c) % 2 != 0:
                    board[r, c] = WHITE_MAN
        return board

    def get_piece(self, r: int, c: int) -> int:
        return self.board[r, c]

    def is_white(self, piece: int) -> bool:
        return piece in [WHITE_MAN, WHITE_KING]

    def is_black(self, piece: int) -> bool:
        return piece in [BLACK_MAN, BLACK_KING]

    def is_king(self, piece: int) -> bool:
        return piece in [WHITE_KING, BLACK_KING]

    def get_legal_moves(self, player: int) -> List[Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]]:
        """
        Returns list of (start_pos, end_pos, captured_positions).
        Capture moves are mandatory. Maximum capture rule applies.
        """
        captures = self._get_all_captures(player)
        if captures:
            # Frisian Rule: Must take the maximum number of pieces.
            max_captures = max(len(m[2]) for m in captures)
            return [m for m in captures if len(m[2]) == max_captures]
        
        return self._get_all_normal_moves(player)

    def _get_all_captures(self, player: int) -> List[Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]]:
        all_jumps = []
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                piece = self.board[r, c]
                if (player == WHITE_MAN and self.is_white(piece)) or \
                   (player == BLACK_MAN and self.is_black(piece)):
                    jumps = self._get_piece_jumps(r, c)
                    all_jumps.extend(jumps)
        return all_jumps

    def _get_piece_jumps(self, r: int, c: int, captured_already: List[Tuple[int, int]] = None) -> List[Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]]:
        if captured_already is None:
            captured_already = []
        
        piece = self.board[r, c]
        jumps = []
        
        # Diagonal directions (standard)
        diag_dirs = [(-1, -1), (-1, 1), (1, -1), (1, 1)]
        # Orthogonal directions (Frisian specialty - distance doubled because of black squares grid)
        ortho_dirs = [(0, 2), (0, -2), (2, 0), (-2, 0)]
        
        is_king = self.is_king(piece)
        
        # Handle all directions
        for is_ortho, directions in [(False, diag_dirs), (True, ortho_dirs)]:
            step_size = 2 if is_ortho else 1 # Step to the next black cell in that direction
            
            for dr, dc in directions:
                if is_king:
                    found_enemy = False
                    enemy_pos = None
                    # Search along the line
                    for i in range(1, BOARD_SIZE):
                        nr, nc = r + dr * i, c + dc * i
                        if not (0 <= nr < BOARD_SIZE and 0 <= nc < BOARD_SIZE):
                            break
                        
                        target = self.board[nr, nc]
                        if target == EMPTY:
                            if found_enemy:
                                # Valid landing square after capture
                                new_captured = captured_already + [enemy_pos]
                                
                                # Simulate
                                original_piece = self.board[r, c]
                                self.board[r, c] = EMPTY
                                self.board[nr, nc] = piece
                                
                                sub_jumps = self._get_piece_jumps(nr, nc, new_captured)
                                if sub_jumps:
                                    for sj in sub_jumps:
                                        jumps.append(((r, c), sj[1], sj[2]))
                                else:
                                    jumps.append(((r, c), (nr, nc), new_captured))
                                
                                # Backtrack
                                self.board[r, c] = original_piece
                                self.board[nr, nc] = EMPTY
                            else:
                                continue
                        elif (self.is_white(piece) and self.is_black(target)) or \
                             (self.is_black(piece) and self.is_white(target)):
                            if found_enemy: # Double enemies in a line or enemy already seen
                                break
                            if (nr, nc) in captured_already:
                                break
                            found_enemy = True
                            enemy_pos = (nr, nc)
                        else: # Own piece or edge
                            break
                else:
                    # NORMAL MAN JUMP
                    mid_r, mid_c = r + dr, c + dc
                    end_r, end_c = r + dr * 2, c + dc * 2 # Landing
                    
                    if 0 <= end_r < BOARD_SIZE and 0 <= end_c < BOARD_SIZE:
                        mid_target = self.board[mid_r, mid_c]
                        end_target = self.board[end_r, end_c]
                        
                        if end_target == EMPTY and mid_target != EMPTY and \
                           ((self.is_white(piece) and self.is_black(mid_target)) or \
                            (self.is_black(piece) and self.is_white(mid_target))) and \
                           (mid_r, mid_c) not in captured_already:
                            
                            new_captured = captured_already + [(mid_r, mid_c)]
                            # Simulate
                            orig_start = self.board[r, c]
                            self.board[r, c] = EMPTY
                            self.board[end_r, end_c] = piece
                            
                            sub_jumps = self._get_piece_jumps(end_r, end_c, new_captured)
                            if sub_jumps:
                                for sj in sub_jumps:
                                    jumps.append(((r, c), sj[1], sj[2]))
                            else:
                                jumps.append(((r, c), (end_r, end_c), new_captured))
                            
                            # Backtrack
                            self.board[r, c] = orig_start
                            self.board[end_r, end_c] = EMPTY
        return jumps

    def _get_all_normal_moves(self, player: int) -> List[Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]]:
        moves = []
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                piece = self.board[r, c]
                if (player == WHITE_MAN and self.is_white(piece)) or \
                   (player == BLACK_MAN and self.is_black(piece)):
                    moves.extend(self._get_piece_normal_moves(r, c))
        return moves

    def _get_piece_normal_moves(self, r: int, c: int) -> List[Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]]:
        piece = self.board[r, c]
        moves = []
        
        if self.is_king(piece):
            directions = [(-1, -1), (-1, 1), (1, -1), (1, 1)] # Normal moves only diagonal? Or orthogonal too?
            # Standard Frisian: King moves any distance DIAGONALLY only. Orthogonal is ONLY for captures.
            # King moves diagonally only. Orthogonal is only for captures.
            for dr, dc in directions:
                for dist in range(1, BOARD_SIZE):
                    nr, nc = r + dr * dist, c + dc * dist
                    if 0 <= nr < BOARD_SIZE and 0 <= nc < BOARD_SIZE and self.board[nr, nc] == EMPTY:
                        moves.append(((r, c), (nr, nc), []))
                    else:
                        break
        else:
            # Normal piece moves forward diagonal
            direction_r = -1 if self.is_white(piece) else 1
            for dc in [-1, 1]:
                nr, nc = r + direction_r, c + dc
                if 0 <= nr < BOARD_SIZE and 0 <= nc < BOARD_SIZE and self.board[nr, nc] == EMPTY:
                    moves.append(((r, c), (nr, nc), []))
        return moves

    def make_move(self, move: Tuple[Tuple[int, int], Tuple[int, int], List[Tuple[int, int]]]):
        (start_r, start_c), (end_r, end_c), captured = move
        piece = self.board[start_r, start_c]
        
        self.board[start_r, start_c] = EMPTY
        for cr, cc in captured:
            self.board[cr, cc] = EMPTY
        
        # Promotion
        if not self.is_king(piece):
            if (self.is_white(piece) and end_r == 0) or (self.is_black(piece) and end_r == BOARD_SIZE - 1):
                piece = WHITE_KING if self.is_white(piece) else BLACK_KING
                
        self.board[end_r, end_c] = piece
        
        # Switch turn
        self.turn = BLACK_MAN if self.turn == WHITE_MAN else WHITE_MAN
        
        # Check winner state (Giveaway variant)
        self._check_winner()

    def _check_winner(self):
        # A player wins if they have NO legal moves (all pieces captured or blocked)
        white_moves = self.get_legal_moves(WHITE_MAN)
        black_moves = self.get_legal_moves(BLACK_MAN)
        
        # In Poddavki, having no moves is a WIN. 
        # Check for block (no legal moves for current player)
        if self.turn == WHITE_MAN and not white_moves:
            self.winner = "White (Win by block/no pieces)"
        elif self.turn == BLACK_MAN and not black_moves:
            self.winner = "Black (Win by block/no pieces)"
        
        # Also check if no pieces are left for the OTHER player (immediate win)
        # This prevents the "extra step" issue.
        white_pieces = sum(1 for r in range(BOARD_SIZE) for c in range(BOARD_SIZE) if self.is_white(self.board[r,c]))
        black_pieces = sum(1 for r in range(BOARD_SIZE) for c in range(BOARD_SIZE) if self.is_black(self.board[r,c]))
        
        if white_pieces == 0:
            self.winner = "White (Win by clearing all pieces)"
        elif black_pieces == 0:
            self.winner = "Black (Win by clearing all pieces)"

    def get_state(self) -> np.ndarray:
        return self.board.copy()
