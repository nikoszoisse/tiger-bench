# TimescaleDB Benchmark tool

- Name: Nikolaos Zois
- Total Duration: ~ 12h
- Experience with CLI: Mid-level

## Getting started

Make sure that you already have an updated [**docker**](https://docs.docker.com/engine/install/ubuntu/) or the [**docker-compose**](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-20-04#step-1-installing-docker-compose) command.

Run the `docker compose up` or `docker-compose up` to run the services

The service container which runs the benchmark is named `app`

**What you expect** to see is, **STATS** avg, median, std, min, max, and total time.

----
### Install
```shell
go get github.com/nikoszoisse/tiger-bench &&
go install github.com/nikoszoisse/tiger-bench
```

### Client Usage Example
```shell
tiger-bench -db localhost:5432,homework,interview_user,123 < ./scripts/query_params.csv
```
OR
```shell
tiger-bench -db localhost:5432,homework,interview_user,123 -file ./scripts/query_params.csv
````

### CLI Options

| Option   | Description                                              | Default      |
|----------|----------------------------------------------------------|--------------|
| -v       | Shows debug purpose logs                                 | false        |
| -db      | DB Url format: `server_host:port,database,user,password` | **required** |
| -file    | path to the query_params`.csv` file                      | stdin        |
| -workers | number of concurrent workers executes queries            | num of cpus  |

### Query Params sample file
```csv
hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
```
---
## Dev
### Project packages

- cmd
  - the cli
- internal
  - config, metrics service, metric handlers, csv parser, worker definition
- pkg
  - (helper packages such as logger)

### Design key points
#### Parser
* Lazily read the file
* Process file in pages (dynamically calculated page-size based on machine's memory and file's data)

#### Metrics Service
* Handlers design. Each metric( avg, min, etc.) is a handler with its state
* Every metric record is processed from all metric handlers
* Metric service may buffer as much as `page-size` calculated previously
* Standard Deviation metric (Show how the data are spread around avg)


#### Workers
* Use of a consistent hash algorithm to route queries based on the hostname
  * Prevent the entire load from being skewed on some workers which can cause the system to wear out or stop working.
    * Multiply the routes by log(NumOfWorkers)

### Docker-Setup
* Introduce `docker-compose` file
  * Includes `app`, `pg_client`, and `timescaledb` service
  * Run the `pg_client` once `timescaledb` is _healthy_
  * Run the `app` once `pg_client` is done
