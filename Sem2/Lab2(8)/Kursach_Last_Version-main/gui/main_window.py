import sys
import os

# Add the project root to sys.path
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from ai.agent import RLAgent
from PyQt6.QtWidgets import QApplication, QMainWindow, QWidget, QGridLayout, QPushButton, QVBoxLayout, QLabel, QComboBox, QHBoxLayout, QMessageBox, QFrame, QCheckBox
from PyQt6.QtCore import Qt, QTimer, QSize, QPropertyAnimation, QEasingCurve, QPoint
from PyQt6.QtGui import QPixmap, QIcon, QPainter, QColor, QPen, QRadialGradient, QLinearGradient
from core.engine import FrisianCheckersLogic, WHITE_MAN, BLACK_MAN, BOARD_SIZE
from scripts.benchmark import RandomBot

class BoardCell(QPushButton):
    def __init__(self, r, c):
        super().__init__()
        self.r = r
        self.c = c
        self.setFixedSize(72, 72)
        self.piece = 0
        self.highlight = None
        self.theme_data = None
        self.piece_style_idx = 0
        self.update_style()

    def paintEvent(self, event):
        super().paintEvent(event)
        if not self.theme_data:
            return

        painter = QPainter(self)
        painter.setRenderHint(QPainter.RenderHint.Antialiasing)
        
        # Draw Circular Highlights
        if self.highlight in ["selected", "possible", "mandatory", "captured"]:
            if self.highlight == "selected":
                color = QColor(self.theme_data['highlight_sel'])
                pen = QPen(color, 4)
            elif self.highlight == "possible":
                color = QColor(self.theme_data['highlight_pos'])
                pen = QPen(color, 3, Qt.PenStyle.DashLine)
            elif self.highlight == "captured":
                color = QColor("#e74c3c") # Red for captured
                pen = QPen(color, 5)
            else:
                color = QColor(self.theme_data['highlight_mand'])
                pen = QPen(color, 4)
                
            painter.setPen(pen)
            painter.setBrush(Qt.BrushStyle.NoBrush)
            painter.drawEllipse(self.rect().adjusted(3, 3, -3, -3))
            
            if self.highlight == "captured":
                overlay = QColor(231, 76, 60, 100) # Semi-transparent red
                painter.setBrush(overlay)
                painter.drawEllipse(self.rect().adjusted(5, 5, -5, -5))

        if self.piece == 0:
            return
            
        if self.piece_style_idx == 1:
            # Let the button's icon rendering handle the PNG
            return
        
        is_white = self.piece in [1, 3]
        is_king = self.piece in [3, 4]
        theme_name = self.theme_data.get('name', 'Classic')
        
        rect = self.rect().adjusted(6, 6, -6, -6)
        cx, cy = rect.center().x(), rect.center().y()
        radius = rect.width() / 2
        
        if theme_name == 'Cyberpunk':
            color = QColor("#00fff2") if is_white else QColor("#ff00ff")
            
            # Layer 1: Ambient Glow
            glow = QRadialGradient(cx, cy, radius * 1.5)
            glow.setColorAt(0, QColor(color.red(), color.green(), color.blue(), 100))
            glow.setColorAt(1, Qt.GlobalColor.transparent)
            painter.setBrush(glow)
            painter.setPen(Qt.PenStyle.NoPen)
            painter.drawEllipse(self.rect())

            # Layer 2: Metallic Base
            base_grad = QRadialGradient(cx, cy, radius)
            base_grad.setColorAt(0, QColor("#1a1a2e"))
            base_grad.setColorAt(1, QColor("#050510"))
            painter.setBrush(base_grad)
            painter.setPen(QPen(color, 2))
            painter.drawEllipse(rect)
            
            # Tech Pattern
            painter.setPen(QPen(color.darker(120), 1))
            for i in range(1, 4):
                painter.drawEllipse(rect.adjusted(i*5, i*5, -i*5, -i*5))
            
            # Center Core
            painter.setBrush(color)
            painter.setPen(Qt.PenStyle.NoPen)
            painter.drawEllipse(cx-4, cy-4, 8, 8)

            if is_king:
                painter.setPen(QPen(color.lighter(150), 3))
                painter.drawLine(cx-12, cy, cx+12, cy)
                painter.drawLine(cx, cy-12, cx, cy+12)
        
        else:
            # Realistic Piece Rendering
            # Drop Shadow
            painter.setBrush(QColor(0, 0, 0, 100))
            painter.setPen(Qt.PenStyle.NoPen)
            painter.drawEllipse(rect.translated(3, 3))

            # Material Gradient (Multi-stop)
            grad = QRadialGradient(cx - radius/4, cy - radius/4, radius * 1.5)
            if is_white:
                grad.setColorAt(0, QColor("#ffffff"))   # Highlight
                grad.setColorAt(0.3, QColor("#f0f0f0")) # Mid
                grad.setColorAt(0.7, QColor("#d0d0d0")) # Shadow
                grad.setColorAt(1, QColor("#a0a0a0"))   # Edge
                rim = QColor("#95a5a6")
            else:
                grad.setColorAt(0, QColor("#555555"))   # Highlight
                grad.setColorAt(0.3, QColor("#222222")) # Body
                grad.setColorAt(1, QColor("#000000"))   # Deep shadow
                rim = QColor("#2c3e50")
            
            painter.setBrush(grad)
            painter.setPen(QPen(rim, 1))
            painter.drawEllipse(rect)
            
            # Surface Texture
            if is_white:
                painter.setPen(QPen(QColor(0, 0, 0, 40), 1))
            else:
                painter.setPen(QPen(QColor(255, 255, 255, 30), 1))
            for i in range(3, 8):
                painter.drawEllipse(rect.adjusted(i*2, i*2, -i*2, -i*2))
            
            # Glossy Highlight
            glinf = QLinearGradient(cx, cy - radius, cx, cy)
            glinf.setColorAt(0, QColor(255, 255, 255, 100))
            glinf.setColorAt(0.5, QColor(255, 255, 255, 0))
            painter.setBrush(glinf)
            painter.setPen(Qt.PenStyle.NoPen)
            painter.drawEllipse(rect.adjusted(5, 2, -5, int(radius * 0.8)))

            if is_king:
                painter.setBrush(Qt.BrushStyle.NoBrush)
                # Outer gold ring
                painter.setPen(QPen(QColor("#f1c40f"), 4))
                painter.drawEllipse(cx-14, cy-14, 28, 28)
                # Inner gold glow
                painter.setPen(QPen(QColor("#f39c12"), 2))
                painter.drawEllipse(cx-10, cy-10, 20, 20)

    def update_style(self, piece=0, highlight=None, theme=None, style_idx=0):
        self.piece = piece
        self.highlight = highlight
        self.theme_data = theme
        self.piece_style_idx = style_idx
        
        if not theme: return
        
        is_dark = (self.r + self.c) % 2 != 0
        base_color = theme['board_dark'] if is_dark else theme['board_light']
        
        style = f"background-color: {base_color}; border: none; border-radius: 0px;"
        self.setStyleSheet(style)
        
        if self.piece_style_idx == 1 and self.piece != 0 and 'assets' in theme:
            px_path = theme['assets'].get(self.piece)
            if px_path and os.path.exists(px_path):
                self.setIcon(QIcon(px_path))
                self.setIconSize(QSize(55, 55))
            else:
                self.setIcon(QIcon())
        else:
            self.setIcon(QIcon())
            
        self.update() # Trigger paintEvent

class CheckersWindow(QMainWindow):
    THEMES = {
        'Classic': {
            'name': 'Classic',
            'bg': '#1a1a1a',
            'board_light': '#d2b48c',
            'board_dark': '#8b5a2b',
            'highlight_sel': '#f1c40f',
            'highlight_pos': '#2ecc71',
            'highlight_mand': '#e74c3c',
            'border': '#5d4037',
            'assets': {1: "assets/white_checkers.png", 2: "assets/black_checkers.png", 3: "assets/white_checkers_king.png", 4: "assets/black_checkers_king.png"}
        },
        'Cyberpunk': {
            'name': 'Cyberpunk',
            'bg': '#050505',
            'board_light': '#1a1a2e',
            'board_dark': '#16213e',
            'highlight_sel': '#00fff2',
            'highlight_pos': '#0f3460',
            'highlight_mand': '#e94560',
            'border': '#00fff2',
            'assets': {1: "assets/cyber_w.png", 2: "assets/cyber_b.png", 3: "assets/cyber_wk.png", 4: "assets/cyber_bk.png"}
        },
        'Emerald': {
            'name': 'Emerald',
            'bg': '#0a1a0a',
            'board_light': '#2d5a27',
            'board_dark': '#1e3d1a',
            'highlight_sel': '#58d68d',
            'highlight_pos': '#239b56',
            'highlight_mand': '#ec7063',
            'border': '#2d5a27',
            'assets': {1: "assets/white_checkers.png", 2: "assets/black_checkers.png", 3: "assets/white_checkers_king.png", 4: "assets/black_checkers_king.png"}
        }
    }

    def __init__(self):
        super().__init__()
        self.setWindowTitle("Фризские шашки")
        self.engine = FrisianCheckersLogic()
        self.selected_pos = None
        self.cells = {}
        self.current_theme = 'Classic'
        self.current_style_idx = 0 # 0=Vector, 1=PNG
        
        # UI Components and state
        self.lang = 'RU'
        self.translations = {
            'RU': {
                'title': "Фризские шашки",
                'mode': "Режим:",
                'lang': "Язык:",
                'reset': "Сброс игры",
                'white_turn': "Ход белых",
                'black_turn': "Ход черных",
                'game_over': "ИГРА ОКОНЧЕНА!",
                'white_win': "Белые (Победа по блоку/без фигур)",
                'black_win': "Черные (Победа по блоку/без фигур)",
                'msg_title': "Конец игры",
                'msg_text': "Игра окончена!\n\nПобедитель: {winner}",
                'modes': [
                    "Человек (Белые) vs Человек (Черные)", 
                    "Человек (Белые) vs Random AI (Черные)", 
                    "Человек (Белые) vs RL AI (Черные)", 
                    "Random AI (Белые) vs Random AI (Черные)",
                    "Random AI (Белые) vs RL AI (Черные)",
                    "RL AI (Белые) vs Random AI (Черные)",
                    "RL AI (Белые) vs RL AI (Черные)"
                ],
                'themes': ["Классическая", "Киберпанк", "Изумрудная"],
                'theme_lbl': "Тема:",
                'style_lbl': "Стиль шашек:",
                'anim_lbl': "Анимации",
                'styles': ["Вектор 4K", "Классика (PNG)"]
            },
            'EN': {
                'title': "Frisian Checkers",
                'mode': "Mode:",
                'lang': "Lang:",
                'reset': "Reset Game",
                'white_turn': "White's Turn",
                'black_turn': "Black's Turn",
                'game_over': "GAME OVER!",
                'white_win': "White (Win by block/no pieces)",
                'black_win': "Black (Win by block/no pieces)",
                'msg_title': "Game Over",
                'msg_text': "The game is over!\n\nWinner: {winner}",
                'modes': [
                    "Human (White) vs Human (Black)", 
                    "Human (White) vs Random AI (Black)", 
                    "Human (White) vs RL AI (Black)", 
                    "Random AI (White) vs Random AI (Black)",
                    "Random AI (White) vs RL AI (Black)",
                    "RL AI (White) vs Random AI (Black)",
                    "RL AI (White) vs RL AI (Black)"
                ],
                'themes': ["Classic", "Cyberpunk", "Emerald"],
                'theme_lbl': "Theme:",
                'style_lbl': "Piece Style:",
                'anim_lbl': "Animations",
                'styles': ["Vector 4K", "Classic (PNG)"]
            }
        }
        
        self.is_animating = False
        self.pending_move = None
        
        # Bots
        self.random_bot = RandomBot()
        self.rl_agent = RLAgent()
        
        self.init_ui()
        self.update_ui_text()
        
        # Timer for AI moves
        self.ai_timer = QTimer()
        self.ai_timer.timeout.connect(self.check_ai_turn)
        self.ai_timer.start(500)

    def update_ui_text(self):
        t = self.translations[self.lang]
        self.setWindowTitle(t['title'])
        self.mode_label.setText(t['mode'])
        self.lang_label.setText(t['lang'])
        self.theme_label.setText(t['theme_lbl'])
        self.style_label.setText(t['style_lbl'])
        self.anim_checkbox.setText(t['anim_lbl'])
        self.reset_btn.setText(t['reset'])
        self.refresh_board()
        
        # Update mode items
        self.mode_selector.blockSignals(True)
        current_mode = self.mode_selector.currentIndex()
        self.mode_selector.clear()
        self.mode_selector.addItems(t['modes'])
        self.mode_selector.setCurrentIndex(max(0, current_mode))
        self.mode_selector.blockSignals(False)

        # Update theme items
        self.theme_selector.blockSignals(True)
        current_theme = self.theme_selector.currentIndex()
        self.theme_selector.clear()
        self.theme_selector.addItems(t['themes'])
        self.theme_selector.setCurrentIndex(max(0, current_theme))
        self.theme_selector.blockSignals(False)

        # Update style items
        self.style_selector.blockSignals(True)
        current_style = self.style_selector.currentIndex()
        self.style_selector.clear()
        self.style_selector.addItems(t['styles'])
        self.style_selector.setCurrentIndex(max(0, current_style))
        self.style_selector.blockSignals(False)

    def change_language(self, index):
        self.lang = 'RU' if index == 0 else 'EN'
        self.update_ui_text()

    def change_theme(self, index):
        theme_names = ['Classic', 'Cyberpunk', 'Emerald']
        self.current_theme = theme_names[index]
        theme_data = self.THEMES[self.current_theme]
        
        # Update main window and panel styles
        bg_col = theme_data['bg']
        self.centralWidget().setStyleSheet(f"background-color: {bg_col};")
        
        # Update info label style (Themes use highlight_sel as a signature color)
        accent = theme_data['highlight_sel']
        self.info_label.setStyleSheet(f"font-weight: bold; font-size: 28px; margin: 20px; color: {accent}; letter-spacing: 2px;")
        
        # Update board frame
        self.board_frame.setStyleSheet(f"border: 4px solid {theme_data['border']}; border-radius: 10px; padding: 5px; background: {theme_data['board_dark']};")
        
        # Update stat labels
        stat_style = f"font-size: 20px; font-weight: bold; color: #ffffff; background: {theme_data['board_dark']}; border: 1px solid {accent}; padding: 15px; border-radius: 10px; min-width: 170px;"
        self.white_score.setStyleSheet(stat_style)
        self.black_score.setStyleSheet(stat_style)
        
        # Update reset button accent
        self.reset_btn.setStyleSheet(f"background: #c0392b; color: white; font-weight: bold; border-radius: 5px; padding: 10px 20px; font-size: 16px; border: 1px solid {accent};")
        
        self.refresh_board()

    def change_piece_style(self, index):
        self.current_style_idx = index
        self.refresh_board()

    def init_ui(self):
        central = QWidget()
        central.setStyleSheet("background-color: #1a1a1a;")
        self.setCentralWidget(central)
        layout = QVBoxLayout(central)
        layout.setContentsMargins(30, 30, 30, 30)
        
        # Global Styles
        self.setStyleSheet("""
            QMainWindow { background-color: #1a1a1a; }
            QLabel { color: #ecf0f1; font-family: 'Segoe UI', Arial; }
            QComboBox { 
                background: #34495e; color: white; padding: 5px; border-radius: 3px; 
                min-width: 150px; border: 1px solid #5d6d7e;
            }
            QPushButton#reset_btn { 
                background: #c0392b; color: white; font-weight: bold; border-radius: 3px; 
                padding: 5px 15px; 
            }
            QPushButton#reset_btn:hover { background: #e74c3c; }
        """)
        
        # Controls
        ctrl_layout = QHBoxLayout()
        
        # Mode
        self.mode_label = QLabel("Mode:")
        self.mode_selector = QComboBox()
        # Items are added in update_ui_text()
        
        # Language
        self.lang_label = QLabel("Language:")
        self.lang_selector = QComboBox()
        self.lang_selector.addItems(["RU", "EN"])
        self.lang_selector.currentIndexChanged.connect(self.change_language)
        
        # Theme
        self.theme_label = QLabel("Theme:")
        self.theme_selector = QComboBox()
        # Items added in update_ui_text
        self.theme_selector.currentIndexChanged.connect(self.change_theme)
        
        # Piece Style
        self.style_label = QLabel("Piece Style:")
        self.style_selector = QComboBox()
        self.style_selector.currentIndexChanged.connect(self.change_piece_style)
        
        # Animations
        self.anim_checkbox = QCheckBox("Animations")
        self.anim_checkbox.setChecked(True)
        self.anim_checkbox.setStyleSheet("color: white; font-weight: bold;")
        
        ctrl_layout.addStretch()
        ctrl_layout.addWidget(self.mode_label)
        ctrl_layout.addWidget(self.mode_selector)
        ctrl_layout.addSpacing(15)
        ctrl_layout.addWidget(self.lang_label)
        ctrl_layout.addWidget(self.lang_selector)
        ctrl_layout.addSpacing(15)
        ctrl_layout.addWidget(self.theme_label)
        ctrl_layout.addWidget(self.theme_selector)
        ctrl_layout.addSpacing(15)
        ctrl_layout.addWidget(self.style_label)
        ctrl_layout.addWidget(self.style_selector)
        ctrl_layout.addSpacing(15)
        ctrl_layout.addWidget(self.anim_checkbox)
        ctrl_layout.addSpacing(30)
        
        self.reset_btn = QPushButton("Reset Game")
        self.reset_btn.setObjectName("reset_btn")
        self.reset_btn.clicked.connect(self.reset_game)
        self.reset_btn.setMinimumHeight(40)
        ctrl_layout.addWidget(self.reset_btn)
        ctrl_layout.addStretch()
        
        layout.addLayout(ctrl_layout)
        
        self.info_label = QLabel("White's Turn")
        self.info_label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.info_label.setStyleSheet("font-weight: bold; font-size: 24px; margin: 20px; color: #f1c40f;")
        layout.addWidget(self.info_label)
        
        grid = QGridLayout()
        grid.setSpacing(0)
        grid.setAlignment(Qt.AlignmentFlag.AlignCenter)
        
        # Main content area with side stats
        main_content = QHBoxLayout()
        
        # Left Panel (Stats)
        left_panel = QVBoxLayout()
        left_panel.setContentsMargins(20, 0, 20, 0)
        
        self.white_score = QLabel("White: 20")
        self.black_score = QLabel("Black: 20")
        for lbl in [self.white_score, self.black_score]:
            lbl.setStyleSheet("font-size: 20px; font-weight: bold; color: #ecf0f1; background: #2c3e50; padding: 15px; border-radius: 8px; min-width: 150px;")
            lbl.setAlignment(Qt.AlignmentFlag.AlignCenter)
            left_panel.addWidget(lbl)
            
        left_panel.addStretch()
        
        # Board Frame (Decorative)
        self.board_frame = QFrame()
        self.board_frame.setObjectName("board_frame")
        self.board_frame.setLayout(grid)
        self.board_frame.setLineWidth(5)
        
        main_content.addStretch()
        main_content.addLayout(left_panel)
        main_content.addSpacing(40)
        main_content.addWidget(self.board_frame)
        main_content.addStretch()
        
        # Create cells
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                cell = BoardCell(r, c)
                cell.clicked.connect(lambda _, r=r, c=c: self.on_cell_click(r, c))
                grid.addWidget(cell, r, c)
                self.cells[(r, c)] = cell
        
        layout.addLayout(main_content)
        self.refresh_board()

    def reset_game(self):
        self.engine = FrisianCheckersLogic()
        self.selected_pos = None
        self.refresh_board()

    def check_ai_turn(self):
        if self.engine.winner or self.is_animating: return
        
        mode_idx = self.mode_selector.currentIndex()
        turn = self.engine.turn
        
        # Game logic control
        delay = 600 if self.anim_checkbox.isChecked() else 100
        
        if mode_idx == 1 and turn == BLACK_MAN:
            self.make_bot_move(self.random_bot)
        elif mode_idx == 2 and turn == BLACK_MAN:
            self.make_bot_move(self.rl_agent)
        elif mode_idx == 3:
            QTimer.singleShot(delay, lambda: self.make_bot_move(self.random_bot))
        elif mode_idx == 4:
            if turn == WHITE_MAN: QTimer.singleShot(delay, lambda: self.make_bot_move(self.random_bot))
            else: QTimer.singleShot(delay, lambda: self.make_bot_move(self.rl_agent))
        elif mode_idx == 5:
            if turn == WHITE_MAN: QTimer.singleShot(delay, lambda: self.make_bot_move(self.rl_agent))
            else: QTimer.singleShot(delay, lambda: self.make_bot_move(self.random_bot))
        elif mode_idx == 6:
            QTimer.singleShot(delay, lambda: self.make_bot_move(self.rl_agent))

    def make_bot_move(self, bot):
        if self.engine.winner or self.is_animating: return
        move = bot.get_move(self.engine, self.engine.turn)
        if move:
            # Highlight start cell for "thinking" effect if RL bot
            is_rl = hasattr(bot, 'model') # Simple check for RL agent
            if is_rl and self.anim_checkbox.isChecked():
                start_pos = move[0]
                self.cells[start_pos].update_style(
                    self.engine.board[start_pos],
                    "selected",
                    self.THEMES[self.current_theme],
                    self.current_style_idx
                )
                QTimer.singleShot(200, lambda: self.animate_move(move))
            else:
                self.animate_move(move)

    def animate_move(self, move):
        if not self.anim_checkbox.isChecked():
            self.engine.make_move(move)
            self.refresh_board()
            return

        self.is_animating = True
        self.pending_move = move
        
        (start_r, start_c), (end_r, end_c), captured = move
        
        start_cell = self.cells[(start_r, start_c)]
        end_cell = self.cells[(end_r, end_c)]
        
        # Get absolute positions
        start_pos = start_cell.mapTo(self.board_frame, QPoint(0, 0))
        end_pos = end_cell.mapTo(self.board_frame, QPoint(0, 0))
        
        # Create a flying visual mockup
        self.anim_piece = BoardCell(0, 0)
        self.anim_piece.setParent(self.board_frame)
        self.anim_piece.setFixedSize(start_cell.size())
        
        # Copy style exactly
        piece_val = self.engine.board[start_r, start_c]
        theme_data = self.THEMES[self.current_theme]
        self.anim_piece.update_style(piece_val, None, theme_data, self.current_style_idx)
        
        # Make transparent background manually since it inherits the board cell bg
        self.anim_piece.setStyleSheet("background-color: transparent; border: none;")
        
        # Hide original piece temporarily
        start_cell.update_style(0, None, theme_data, self.current_style_idx)
        
        self.anim_piece.move(start_pos)
        self.anim_piece.show()
        self.anim_piece.raise_()
        
        # Highlight captured pieces immediately
        for cr, cc in captured:
            self.cells[(cr, cc)].update_style(
                self.engine.board[cr, cc], 
                "captured", 
                theme_data, 
                self.current_style_idx
            )
        
        self.anim = QPropertyAnimation(self.anim_piece, b"pos")
        self.anim.setDuration(400) # Slightly slower for clarity
        self.anim.setStartValue(start_pos)
        self.anim.setEndValue(end_pos)
        self.anim.setEasingCurve(QEasingCurve.Type.InOutCubic)
        self.anim.finished.connect(self.on_animation_finished)
        self.anim.start()

    def on_animation_finished(self):
        # Process captures before removal for visual transition
        captured_positions = self.pending_move[2]
        
        if captured_positions:
            # Short delay to see the result of the jump before pieces vanish
            QTimer.singleShot(200, self.finalize_move)
        else:
            self.finalize_move()

    def finalize_move(self):
        if hasattr(self, 'anim_piece'):
            self.anim_piece.deleteLater()
            
        self.engine.make_move(self.pending_move)
        self.pending_move = None
        self.is_animating = False
        self.refresh_board()

    def refresh_board(self):
        # Identify possible moves if a piece is selected
        possible_targets = []
        legal_moves = self.engine.get_legal_moves(self.engine.turn)
        
        if self.selected_pos:
            possible_targets = [m[1] for m in legal_moves if m[0] == self.selected_pos]
            
        # Identify pieces that MUST capture
        mandatory_sources = []
        is_capture_available = any(len(m) > 2 for m in legal_moves) # Jump move check
        if is_capture_available:
            mandatory_sources = list(set(m[0] for m in legal_moves if len(m) > 2))

        # Validate human turn
        mode_idx = self.mode_selector.currentIndex()
        is_human_turn = False
        if mode_idx == 0: is_human_turn = True
        elif (mode_idx == 1 or mode_idx == 2) and self.engine.turn == WHITE_MAN:
            is_human_turn = True

        white_count = 0
        black_count = 0

        theme_data = self.THEMES[self.current_theme]
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                piece = self.engine.get_piece(r, c)
                pos = (r, c)
                
                # Count pieces
                if self.engine.is_white(piece): white_count += 1
                elif self.engine.is_black(piece): black_count += 1

                highlight = None
                if pos == self.selected_pos:
                    highlight = "selected"
                elif pos in possible_targets:
                    highlight = "possible"
                elif is_human_turn and pos in mandatory_sources:
                    highlight = "mandatory"
                    
                self.cells[pos].update_style(piece, highlight, theme_data, self.current_style_idx)
        
        # Update labels
        self.white_score.setText(f"{'White' if self.lang == 'EN' else 'Белые'}: {white_count}")
        self.black_score.setText(f"{'Black' if self.lang == 'EN' else 'Черные'}: {black_count}")

        t = self.translations[self.lang]
        if self.engine.turn == WHITE_MAN:
            turn_text = t['white_turn']
        else:
            turn_text = t['black_turn']
            
        if self.engine.winner:
            turn_text = f"{t['game_over']} {self.engine.winner}"
            
            # Map winner string for translation
            winner_str = self.engine.winner
            if "White" in winner_str: winner_str = t['white_win']
            elif "Black" in winner_str: winner_str = t['black_win']
            
            msg = QMessageBox(self)
            msg.setWindowTitle(t['msg_title'])
            msg.setText(t['msg_text'].format(winner=winner_str))
            msg.setIcon(QMessageBox.Icon.Information)
            msg.exec()
            
        self.info_label.setText(turn_text)

        # Validate human interaction availability
        mode_idx = self.mode_selector.currentIndex()
        is_human_turn = False
        if mode_idx == 0: is_human_turn = True
        elif (mode_idx == 1 or mode_idx == 2) and self.engine.turn == WHITE_MAN:
            is_human_turn = True
            
        if not is_human_turn: return

        piece = self.engine.get_piece(r, c)
        
        if self.selected_pos is None:
            # Select own piece
            if (self.engine.turn == WHITE_MAN and self.engine.is_white(piece)) or \
               (self.engine.turn == BLACK_MAN and self.engine.is_black(piece)):
                self.selected_pos = (r, c)
        else:
            # Try to move
            legal_moves = self.engine.get_legal_moves(self.engine.turn)
            move_to_make = None
            for move in legal_moves:
                if move[0] == self.selected_pos and move[1] == (r, c):
                    move_to_make = move
                    break
            
            if move_to_make:
                self.selected_pos = None
                self.animate_move(move_to_make)
            else:
                # Re-select or de-select
                if (self.engine.turn == WHITE_MAN and self.engine.is_white(piece)) or \
                   (self.engine.turn == BLACK_MAN and self.engine.is_black(piece)):
                    self.selected_pos = (r, c)
                else:
                    self.selected_pos = None
                    
        self.refresh_board()

if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = CheckersWindow()
    window.show()
    sys.exit(app.exec())
