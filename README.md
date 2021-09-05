# Wallet
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/polarbit/bluelabs-wallet)



### References
- [Echo Graceful Shutdown](https://echo.labstack.com/cookbook/graceful-shutdown/)



### DUMP
- Inside the container, password is not required
$ docker run --name db -e POSTGRES_PASSWORD=1234 -p 5342:5342 -d postgres
$ create database walletdb;
& create table

pgx
- postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

go run . db -d "postgresql://postgres:1234@localhost" --initdb walletdb
go run . db -d "postgresql://postgres:1234@localhost" --dropdb walletdb

# #List dbs
psql> \l
# List table
psql> \dt
select * from pg_catalog.pg_tables where tablename like 'wallet%';

# Changed db name from default to hede
BL_DB_DATABASE=hede go run . config

# Run integration tests
go test ./... --integration -v 