"""Static server for the demo: like `python -m http.server`, plus an ETag header (our downloader needs one).

Run from sdkweblegacy:  python src/infer/demo/serve.py   ->   http://localhost:8000/src/infer/demo/index.html
"""
import os
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer


class Handler(SimpleHTTPRequestHandler):
    def end_headers(self):
        # Always revalidate: the browser must not answer our HEAD requests from its own cache
        self.send_header("Cache-Control", "no-cache")
        path = self.translate_path(self.path)
        if os.path.isfile(path):
            st = os.stat(path)
            self.send_header("ETag", '"%x-%x"' % (int(st.st_mtime), st.st_size))
        super().end_headers()


ThreadingHTTPServer(("", 8000), Handler).serve_forever()
