DROP DATABASE IF EXISTS homework;
DROP USER IF EXISTS interview_user;
DROP TABLE IF EXISTS cpu_usage;

CREATE DATABASE homework;
\c homework
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE cpu_usage(
  ts    TIMESTAMPTZ,
  host  TEXT,
  usage DOUBLE PRECISION
);

CREATE USER interview_user WITH PASSWORD '123';

SELECT create_hypertable('cpu_usage', 'ts');

GRANT CONNECT ON DATABASE homework TO interview_user;
GRANT USAGE ON SCHEMA public TO interview_user;
GRANT SELECT ON cpu_usage TO interview_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO interview_user;