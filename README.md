## How to start

To run FOOD API for local development you need Docker and Docker-Compose

For local development you need start next services:
- app (Golang app) from Dockerfile
- db (MySQL server)


Next command setup **food-api_dev** (Golang App) container with hot reloading and debugger and Postgres DB server
```bash
$ docker-compose -f ./docker-compose.yml -f ./docker-compose-DEV.yml up -d --build app db
```
Now you can access Food API using URL: http://localhost:3000


If you need only setup Food API for using in local environment (without development), run next command:
```bash
$ docker-compose -f ./docker-compose.yml up -d --build app db
```
Now you can access Food API using URL: http://localhost:5000

All local environment variables you can find in file **.env**

---

## How to create and apply SQL migrations

First you need enter inside container  
```bash
$ docker exec -it food-api bash
```

#### Create new SQL migration
To create new migration, you can use following command:
```bash
# /app/bin/migrate create -ext sql -dir /app/migrations NAME
```

Now you can find generated files in folder **migrations**
- {digits}_NAME.up.sql
- {digits}_NAME.down.sql

#### Apply migrations
To apply all existing migrations, you can use following command:
```bash
# /app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" up
```

If you want to see current DB version, you can use following command:
```bash
# /app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" version
```

## Swagger integration
To generate Swagger docs use special tool.
To update documentation:
First you need enter inside container  
```bash
$ docker exec -it food-api bash
```

```bash
$ /app/bin/swag init -g handler/handler.go
```

Swagger docs you can found using URL [http://localhost:3000/swagger/index.html]
More info to generate or update Swagger docs - you can find [https://github.com/swaggo/gin-swagger]

## Manage dependencies
To manage dependencies uses **vgo**
All dependencies requirements are stored in:
- go.mod and
- go.sum files
