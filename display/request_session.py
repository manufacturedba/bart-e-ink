import wifi
import adafruit_connection_manager
import adafruit_requests

class RequestSession():
    def __init__(self, SSID, PASSWORD):
        wifi.radio.connect(SSID, PASSWORD)
        
        pool = adafruit_connection_manager.get_radio_socketpool(wifi.radio)
        ssl_context = adafruit_connection_manager.get_radio_ssl_context(wifi.radio)
        self.wifi = wifi
        self.requests = adafruit_requests.Session(pool, ssl_context)

    def get(self, url):
        return self.requests.get(url)

    def post(self, url, data):
        return self.requests.post(url, data)

    def put(self, url, data):
        return self.requests.put(url, data)

    def delete(self, url):
        return self.requests.delete(url)