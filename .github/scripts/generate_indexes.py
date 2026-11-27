#!/usr/bin/env python3
"""
Generates static index.html directory listing files for GitHub Pages.
"""

import html
import os
import sys


def format_size(size_bytes: int) -> str:
    if size_bytes >= 1024 * 1024:
        return f"{size_bytes / (1024 * 1024):.1f} MB"
    if size_bytes >= 1024:
        return f"{size_bytes / 1024:.1f} KB"
    return f"{size_bytes} B"


def generate_indexes(staging_dir: str):
    staging_dir = os.path.abspath(staging_dir)
    release_dir = os.path.join(staging_dir, "release")
    if not os.path.exists(release_dir):
        print(f"Warning: {release_dir} does not exist, nothing to index.")
        return

    # 1. Root redirect index.html
    root_index = os.path.join(staging_dir, "index.html")
    with open(root_index, "w", encoding="utf-8") as f:
        f.write("""<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="0; url=release/">
  <title>Redirecting to release/</title>
</head>
<body>
  <p>Redirecting to <a href="release/">release/</a>...</p>
</body>
</html>
""")

    # 2. Walk release_dir and generate index.html in each directory
    count = 0
    for dirpath, dirnames, filenames in os.walk(release_dir):
        # Exclude hidden directories
        dirnames[:] = [d for d in sorted(dirnames) if not d.startswith(".")]

        rel_from_staging = os.path.relpath(dirpath, staging_dir).replace(os.sep, "/")
        parts = [p for p in rel_from_staging.split("/") if p]
        total_parts = len(parts)

        # Build relative breadcrumb navigation
        breadcrumbs = [f'<a href="{"../" * total_parts}">/</a>']
        for i, part in enumerate(parts):
            up = total_parts - 1 - i
            if up == 0:
                breadcrumbs.append(f"<span>{html.escape(part)}</span>")
            else:
                breadcrumbs.append(f'<a href="{"../" * up}">{html.escape(part)}</a>')
        bc_html = " / ".join(breadcrumbs) + " /"

        # Build table rows
        rows = []
        if dirpath != release_dir:
            rows.append('<tr><td>📁 <a href="../">../</a></td><td>-</td></tr>')

        for d in dirnames:
            rows.append(f'<tr><td>📁 <a href="{html.escape(d)}/">{html.escape(d)}/</a></td><td>-</td></tr>')

        visible_files = [f for f in sorted(filenames) if f != "index.html" and not f.startswith(".")]
        for f in visible_files:
            fp = os.path.join(dirpath, f)
            size_str = format_size(os.path.getsize(fp))
            rows.append(f'<tr><td>📄 <a href="{html.escape(f)}">{html.escape(f)}</a></td><td>{size_str}</td></tr>')

        page_html = f"""<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Index of /{html.escape(rel_from_staging)}/</title>
  <style>
    body {{
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, monospace;
      padding: 24px;
      color: #24292f;
      background: #ffffff;
      margin: 0;
    }}
    .container {{
      max-width: 900px;
      margin: 0 auto;
    }}
    h2 {{
      font-size: 1.25rem;
      font-weight: 600;
      margin: 0 0 12px 0;
      word-break: break-all;
    }}
    hr {{
      border: 0;
      border-top: 1px solid #d0d7de;
      margin: 16px 0;
    }}
    table {{
      border-collapse: collapse;
      width: 100%;
      font-size: 14px;
      font-family: monospace;
    }}
    th, td {{
      text-align: left;
      padding: 6px 12px;
    }}
    th {{
      border-bottom: 2px solid #d0d7de;
      color: #57606a;
    }}
    tr:hover {{
      background-color: #f6f8fa;
    }}
    a {{
      color: #0969da;
      text-decoration: none;
    }}
    a:hover {{
      text-decoration: underline;
    }}
    td:last-child, th:last-child {{
      text-align: right;
    }}
  </style>
</head>
<body>
  <div class="container">
    <h2>Index of {bc_html}</h2>
    <hr>
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
        </tr>
      </thead>
      <tbody>
        {"".join(rows)}
      </tbody>
    </table>
    <hr>
  </div>
</body>
</html>
"""
        with open(os.path.join(dirpath, "index.html"), "w", encoding="utf-8") as f:
            f.write(page_html)
        count += 1

    print(f"Generated {count} index.html files in {staging_dir}.")


if __name__ == "__main__":
    target = sys.argv[1] if len(sys.argv) > 1 else "gh-pages-staging"
    generate_indexes(target)
