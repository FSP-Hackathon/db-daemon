
SELECT pg_terminate_backend(%s)
FROM pg_stat_activity
WHERE %s != pg_backend_pid()
AND datname IS NOT NULL
AND leader_pid IS NULL;