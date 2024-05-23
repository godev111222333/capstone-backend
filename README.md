# Capstone's backend

## Prerequisite installation: 
1. Go version >= 1.21.5
2. Docker
3. Make

## Migration steps
1. ``make dev-up`` to run Postgresql database inside Docker
2. ``make migrate-up`` to run migration scripts
3. ``make dev-run`` to start the API server
