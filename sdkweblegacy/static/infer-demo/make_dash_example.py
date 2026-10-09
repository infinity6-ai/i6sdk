"""Builds dash-example.sqlite, a tiny dash for the demo page: python static/infer-demo/make_dash_example.py"""
import os
import sqlite3

path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "dash-example.sqlite")
if os.path.exists(path):
    os.remove(path)

db = sqlite3.connect(path)
db.executescript("""
CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL);
CREATE TABLE funnel (year_month TEXT, event_type TEXT, from_rec TEXT, qt_events INTEGER);
""")

months = ["2024-01", "2024-02", "2024-03"]
events = {"view": 1000, "click": 300, "add_to_cart": 90, "purchase": 25}
for i, month in enumerate(months):
    for event, qt in events.items():
        db.execute("INSERT INTO funnel VALUES (?, ?, 'true', ?)", (month, event, qt * (i + 1) // 3))
        db.execute("INSERT INTO funnel VALUES (?, ?, 'false', ?)", (month, event, qt * 2 * (i + 1) // 3))

# ?1 = the request input: {"dataset": "...", "params": {"min_year_month": "2024-01", "max_year_month": "2024-03"}}
# ret = {"results": {"rows": [{"event_type", "from_rec", "qt_events"}]}}, like the Go func
db.execute("INSERT INTO namedqueries VALUES (?, ?)", ("dash-product-relevance-fashion-dash-funnel", """
WITH arg(j) AS (SELECT ?1)
SELECT json_object('results', json_object('rows', json((
  SELECT json_group_array(json_object('event_type', event_type, 'from_rec', from_rec, 'qt_events', CAST(qt AS TEXT)))
  FROM (
    SELECT event_type, from_rec, sum(qt_events) AS qt
    FROM funnel
    WHERE year_month BETWEEN coalesce(json_extract((SELECT j FROM arg), '$.params.min_year_month'), '0000')
                         AND coalesce(json_extract((SELECT j FROM arg), '$.params.max_year_month'), '9999')
    GROUP BY event_type, from_rec
    ORDER BY qt DESC
  )
)))) AS ret
"""))
db.commit()
db.close()
