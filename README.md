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

## Api methods
* `GET` `/product` - select all products from database
* `GET` `/product/{id}` - select product from database by {id}
* `POST` `/product` - create product with properties passed from json
```
{
    "name": string,
    "description": string,
    "price": int,
    "quantity": int
}
```
* `PATCH` `/product/{id}` - update product by {id} with properties passed from json
```
{
    "name": string,
    "description": string,
    "price": int,
    "quantity": int
}
```
* `DELETE` `/product/{id}` - delete product by {id}

* `GET` `/customer` - select all customers from database
* `GET` `/customer/{id}` - select customer from database by {id}
* `POST` `/customer` - create customer with properties passed from json
```
{
    "first_name": string,
    "last_name": string
}
```
* `PATCH` `/customer/{id}` - update customer by {id} with properties passed from json
```
{
    "first_name": string,
    "last_name": string
}
```
* `DELETE` `/customer/{id}` - delete customer by {id}

* `GET` `/bill` - select all bills from database
* `GET` `/bill/{id}` - select bill from database by {id}
* `POST` `/bill` - create bill with properties passed from json
```
{
    "customer": int,
    "products": [
        {
            "product": int,
            "quantity" int
        }
    ]
}
```
* `PATCH` `/bill/{id}` - update bill by {id} with properties passed from json
```
{
    "customer": int,
    "products": [
        {
            "product": int,
            "quantity" int
        }
    ]
}
```
* `DELETE` `/bill/{id}` - delete bill by {id}

* `GET` `/bill/{id}/product` - select all products related to bill received by {id}
* `POST` `/bill/{id}/product` - add new product to bill received by {id}
* `DELETE` `/bill/{bill_id}/product/{product_id}` - delete product with id {product_id} from bill with id {bill_id}