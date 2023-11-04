SELECT 
  pid, 
  wait_event, 
  now(),
  backend_start
FROM pg_stat_activity;