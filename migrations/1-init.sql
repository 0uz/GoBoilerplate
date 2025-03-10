SELECT 'CREATE DATABASE appdb'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'appdb')\gexec

DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles WHERE rolname = 'app'
   ) THEN
      CREATE ROLE app WITH LOGIN PASSWORD 'StrongPassword!';
   END IF;
END
$do$;

\c appdb

CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION app;

GRANT CONNECT ON DATABASE appdb TO app;
GRANT USAGE ON SCHEMA app TO app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO app;

ALTER DEFAULT PRIVILEGES IN SCHEMA app 
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app;

ALTER ROLE app SET search_path TO app;
