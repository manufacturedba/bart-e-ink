import busio
import board
import digitalio

from adafruit_epd.il0373 import Adafruit_IL0373

DISPLAY_HEIGHT = 128
DISPLAY_WIDTH = 296

class IL0373:
    def __init__(self, rotation=0):
        spi = busio.SPI(board.SCK, MOSI=board.MOSI, MISO=board.MISO)
        ecs = digitalio.DigitalInOut(board.D15)
        dc = digitalio.DigitalInOut(board.D33)
        srcs = digitalio.DigitalInOut(board.D32)

        self._display = Adafruit_IL0373(DISPLAY_HEIGHT, DISPLAY_WIDTH, spi, cs_pin=ecs, dc_pin=dc, sramcs_pin=srcs, rst_pin=None, busy_pin=None)
        
        self._display.rotation = rotation
        self.font_name = "font5x8.bin"
        
    def fill(self, color):
        self._display.fill(color)
        
    def text(self, text, x, y, color, size=1):
        self._display.text(text, x, y, color, size=size, font_name=self.font_name)
    
    def display(self):
        self._display.display()