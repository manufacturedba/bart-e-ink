
import os
from adafruit_epd.epd import Adafruit_EPD
from request_session import RequestSession
from il0373 import IL0373
import time

ssid = os.getenv("CIRCUITPY_WIFI_SSID")
password = os.getenv("CIRCUITPY_WIFI_PASSWORD")
endpoint = os.getenv("TRANSIT_ENDPOINT")
refresh_interval = 180 # Recommended safe interval for display refresh

if ssid is None or password is None:
    raise ValueError("Set WIFI_SSID and WIFI_PASSWORD environment variables")

il0373 = IL0373(rotation=3)
requests = RequestSession(ssid, password)

print("Connected to network:", requests.wifi.radio.ap_info.ssid)

while True:
    try:
        print("Fetching display text from %s" % endpoint)
        with requests.get(endpoint) as response:
            decoded = response.json()
            
            il0373.fill(Adafruit_EPD.WHITE)
            
            padding = 10
            for i, line in enumerate(decoded):
                print("Writing to display: \"%s\"" % line)
                color = Adafruit_EPD.BLACK if i % 2 == 0 else Adafruit_EPD.RED
                il0373.text(line, 10, padding, Adafruit_EPD.RED, size=2)
                
                # Increase padding for every other line
                if i % 2 == 0:
                    padding += 20
                else:
                    padding += 40
            
            il0373.display()
            time.sleep(refresh_interval)
    except Exception as e:
        print("Error fetching data:", e)
        time.sleep(refresh_interval)