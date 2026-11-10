"""A downstream that answers, but slowly: one byte every 0.3 seconds.

Run it:  python3 slow_download.py
"""
import http.server
import threading
import time
import urllib.request


class Trickle(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        body = b'{"id":7}'
        self.send_response(200)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        for b in body:
            self.wfile.write(bytes([b]))
            self.wfile.flush()
            time.sleep(0.3)

    def log_message(self, *args):
        pass


server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), Trickle)
threading.Thread(target=server.serve_forever, daemon=True).start()
url = f"http://127.0.0.1:{server.server_address[1]}/users/7"

start = time.monotonic()
with urllib.request.urlopen(url, timeout=1) as resp:
    resp.read()
print(f"timeout=1, but the call took {time.monotonic() - start:.1f}s")
