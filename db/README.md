# go-sro-db

![version](https://img.shields.io/github/v/tag/ferdoran/go-sro-db?label=version)
![build-status](https://img.shields.io/github/workflow/status/ferdoran/go-sro-db/Build%20and%20Publish%20Docker%20Image)
![last-commit](https://img.shields.io/github/last-commit/ferdoran/go-sro-db)

Docker image for go-sro database. Using [mysql](https://hub.docker.com/_/mysql).

Make sure to set the following environment variables when running this image:

- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_ROOT_PASSWORD`