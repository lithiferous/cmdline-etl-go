# ETL using command line with Golang

## Problem Statement

Create a command line tool to load CSV files to a database (Postgres), ensure it provides the following methods:
- **add <filename.csv>** - upload a raw file from [data](./data/raw/) to common table called **snapshots**
- **list** - show the upload status for each file from [data](./data/raw/) (inserted already?)
- **sync** - forever loop to upload any new entries for [data](./data/raw/) files as they arrive

Therefore, the application should be able to query the database (insert/retrieve records with pre-defined schema) from the command line, and launch event-driven process to track file arrival on folder, and upload them to Postgres.

## Service Design

### Postgres database

Dockerized database to keep isolation and persist credentials for maintainability.
Initialize the setup with the following command:

```sh
docker-compose up
```

### Go command line service (cmd-etl)

The application offers multiple modules to support querying databases from the terminal. More information can be found in the corresponding [readme](./cmd-etl/README.md).

## Project Structure

The implementation is decoupled in the following modules:

```
├── data
│   ├── pg  <-(mount directory for postgres running in docker)
│   └── raw
│       ├── snapshot_20230101.csv
│       └── snapshot_20230603.csv
└── cmd-etl <-(service in go)
```

## Usage Demo

![usage.gif](./resources/usage.gif)
