# Capstone's backend

## Prerequisite installation: 
1. Docker
2. Make
3. Golang migration: https://github.com/golang-migrate/migrate
4. Curl (For testing purpose only)

## Run API on your local machine
Before running those steps, please ask author for `config.yaml` file
1. ``make dev-up`` to run Postgresql database + API server inside Docker
2. ``make migrate-up`` to run migration scripts
3. ``curl --request GET http://localhost:9876/ping`` to check if the server running
