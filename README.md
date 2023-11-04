```bash
docker run --name dead_container \
  -e POSTGRES_PASSWORD=DeadDatabase \
  -e POSTGRES_USER=dead \
  -e POSTGRES_DB=dead_db \
  -p 20000:5432 -d postgres:15
```

```sql
DROP SCHEMA IS EXISTS dead_schema CASCADE;
CREATE SCHEMA IF NOT EXISTS dead_schema;

CREATE TABLE dead_schema.test (
  id SERIAL PRIMARY KEY,
  value TEXT
);
```

```py
def Gen():
  f = open("gen.sql", "w")
  f.write("""DROP SCHEMA IF EXISTS dead_schema CASCADE;
CREATE SCHEMA IF NOT EXISTS dead_schema;
CREATE TABLE dead_schema.test (
  id SERIAL PRIMARY KEY,
  value TEXT
);
INSERT INTO dead_schema.test (value) VALUES""")
  for i in range(1, 10000):
    f.write(f'(\'value{i}\'),\n')
  f.close()

Gen()
```

```sql
WITH slojno1 AS (
  SELECT value
  FROM dead_schema.test
  WHERE value != '123'
), slojno2 AS (
  SELECT value
  FROM dead_schema.test
  WHERE value != '245'
), slojno3 AS (
  SELECT value
  FROM dead_schema.test
  WHERE value > '245'
)
SELECT slojno1.value, slojno3.value, slojno2.value
FROM slojno1
INNER JOIN slojno3 ON slojno3.value != slojno1.value
INNER JOIN slojno2 ON slojno2.value != slojno3.value
WHERE slojno2.value != CONCAT(slojno3.value, slojno1.value);
```

```sql
SELECT 
  wait_event, 
  pid, 
  now(),
  backend_start 
FROM pg_stat_activity;
```