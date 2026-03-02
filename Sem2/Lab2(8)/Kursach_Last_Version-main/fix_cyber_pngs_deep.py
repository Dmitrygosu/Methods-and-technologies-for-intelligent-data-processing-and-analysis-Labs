import sys
import os
from PyQt6.QtWidgets import QApplication
from PyQt6.QtGui import QImage, QPainter, QPainterPath, QColor, QPen
from PyQt6.QtCore import Qt, QRectF

def crop_and_scale(image_path):
    if not os.path.exists(image_path):
        return
        
    img = QImage(image_path)
    if img.isNull():
        return
        
    img = img.convertToFormat(QImage.Format.Format_ARGB32)
    # Crop images to remove edges
    inset = int(img.width() * 0.14)
    rect = img.rect().adjusted(inset, inset, -inset, -inset)
    
    # Crop it down
    cropped = img.copy(rect)
    
    # Create a perfectly transparent circular mask
    out_img = QImage(cropped.size(), QImage.Format.Format_ARGB32)
    out_img.fill(Qt.GlobalColor.transparent)
    
    painter = QPainter(out_img)
    painter.setRenderHint(QPainter.RenderHint.Antialiasing)
    
    path = QPainterPath()
    rect_f = QRectF(0, 0, float(cropped.width()), float(cropped.height()))
    path.addEllipse(rect_f)
    
    painter.setClipPath(path)
    painter.drawImage(0, 0, cropped)
    
    # Thin edge to smooth pixels
    painter.setClipping(False)
    painter.setPen(QPen(QColor(0, 0, 0, 150), 2))
    painter.drawEllipse(rect_f)
    painter.end()
    
    # Scale back up to original size (e.g., 480x480)
    final_img = out_img.scaled(img.size(), Qt.AspectRatioMode.KeepAspectRatio, Qt.TransformationMode.SmoothTransformation)
    
    final_img.save(image_path)
    print(f"Deep cropped and scaled {image_path}")

if __name__ == '__main__':
    app = QApplication(sys.argv)
    assets = ['assets/cyber_w.png', 'assets/cyber_b.png', 'assets/cyber_wk.png', 'assets/cyber_bk.png']
    for a in assets:
        crop_and_scale(a)
    print("Done deep fixing PNGs!")
