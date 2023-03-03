## To create a new migration

### Using the cli
```shell
migrate create -ext sql -dir db/migrations -seq <migration_name>
```

migrate create -ext sql -dir db/migrations -seq add_indexes_store_products
### Running migrations with Makefile already in this project
```shell
make migration_up
```