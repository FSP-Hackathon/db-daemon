
SELECT pg_terminate_backend($1)
FROM pg_stat_activity
WHERE $1 != pg_backend_pid()
AND datname IS NOT NULL
AND leader_pid IS NULL;