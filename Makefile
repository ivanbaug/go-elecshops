#!make
include .env

migration_up:
	migrate -path ./db/migrations/ -database "postgres://${LDBUSER}:${LDBPASS}@${LDBHOST}:${LDBPORT}/${LDBNAME}?sslmode=disable" -verbose up

migration_down:
	migrate -path ./db/migrations/ -database "postgres://${LDBUSER}:${LDBPASS}@${LDBHOST}:${LDBPORT}/${LDBNAME}?sslmode=disable" -verbose down

migration_fix:
	migrate -path ./db/migrations/ -database "postgres://${LDBUSER}:${LDBPASS}@${LDBHOST}:${LDBPORT}/${LDBNAME}?sslmode=disable" force VERSION

.PHONY: migration_up migration_down migration_fix
