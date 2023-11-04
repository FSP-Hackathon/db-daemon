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