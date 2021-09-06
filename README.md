# Wallet
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/polarbit/bluelabs-wallet)



### References
- [Echo graceful shutdown](https://echo.labstack.com/cookbook/graceful-shutdown/)
- [Seperate tests using build tags](https://mickey.dev/posts/go-build-tags-testing/)



### DUMP
- Inside the container, password is not required
$ docker run --name db -e POSTGRES_PASSWORD=1234 -p 5342:5342 -d postgres


# Create and Drop database
go run . db -d "postgresql://postgres:1234@localhost" --initdb walletdb
go run . db -d "postgresql://postgres:1234@localhost" --dropdb walletdb

# #List dbs
psql> \l
# List table
psql> \dt
select * from pg_catalog.pg_tables where tablename like 'wallet%';

# Environment variables
DB_DATABASE="postgresql://postgres:1234@localhost"
LOGLEVEL=fatal 

# Run integration tests
go test ./...  -v -tags integration
Note: if integration tests run; unit tests will not run (using build tags)


### TODO
- Move config.json and docker-compose into
- Use connection pooling; or add to your notes.
- Repo: Other fields can be tested for min-max length, existance etc.
- custom error and check against custom error
- run validation also in service



- repo | service => getTransaction yazalım
- service => get-latest
- api => get-balance
- api => *
