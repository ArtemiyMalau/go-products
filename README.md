# Go Products Microservice
Golang implementation of Products JSON API service

## Startup

### Startup on host machine
1. Rename [suggested .env.example](.env.example) file to .env
2. Configure .env file under your environment
3. Start application using following command in terminal:
`go build | .\products.exe -migratedb -seeddb`

### Startup via Docker
1. Rename [suggested .env.example](.env.example) file to .env
2. **REMOVE** `LISTEN_PORT` variable
3. Add `DOCKER_HOST_APP_PORT` and `DOCKER_HOST_DB_PORT` for youself. These variables are describe ports through which you may access to containers' applications (see [docker compose file](docker-compose.yml))
4. Start application in container using following command
`docker-compose --env-file ./.env up -d`
- If you want fully restart application, removing containers and images, paste following command

  `docker-compose --env-file ./.env down; docker container prune; docker image rm products_app:latest; docker-compose --env-file ./.env up -d`

## Command-line options
* `-migratedb` - execute [structure.sql](dbo/structure.sql) script to initialize database structure
* `-seeddb` - execute [seeder.sql](dbo/seeder.sql) to fill database by default data