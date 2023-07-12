# go-elecshops

Small api to manage a list of electronic shops. Made with go and postgresql. It's meant to store and retrieve data from shops of electronics components that I use from time to time.
It's inspired in [rpilocator.com](https://www.rpilocator.com) but for all products of local stores and not as quick with the updates, the data is retrieved with an automated script that runs once daily.

## Database
Migration files and instructions are located in `db/` of this project. [Go migrate](https://github.com/golang-migrate/migrate) has to be installed, can be used with the `Makefile` in the root folder to set the db quickly via cli.

## Docker

The project is meant to be run in a debian docker container.

First the project needs to be built with the help of our docker compose file.
```bash
docker compose build
```

After the image is built, we can run lift the docker container with the following command.
```bash
docker compose up -d
```
## Api documentation
Since the project is so small I didn't properly document the api, but I made tests to the endpoints with postman and covered most of the usage cases.
It can be found in https://documenter.getpostman.com/view/13923274/2s946cfZHN .