# Wallet
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/polarbit/bluelabs-wallet)


## How to Run
- You need go:1.17, docker and docker-compose
- First run the database server     `$ docker-compose up -d`
- Create the database and tables    `$ go run . db --initdb`
- Run unit tests                    `$ go test ./... -v`
- Run integration tests             `$ go test ./... -v -tags integration`
- Run the api                       `$ go run . api`   
- Repository can also be opened and run on GitPod (free, but login required).

## Enpoints

```json
# POST /wallets
{
    "externalId" : "userid-123",     // should be unique
    "labels" : {                    // string dictionary for metadata
        "somekey" : "somevalue"
    }
}
// Idempotency is achieved by "externalid" value
```

```json
# GET  /wallets/:id
{
    "id": 9180,
    "labels" : {
        "somekey" : "somevalue"
    },
    "externalId" : "userid-123",
    "created" : "2021-09-07T01:42:00"
}
``` 

```json
# POST  /wallets/:id/transactions  
{
    "amount" : 10.0,             // amount may be negative or positive,
    "fingerprint" : "TX123A001", // like idempotency key; should be unique
    "labels" : {                 // string dictionary for metadata
        "couponId" : "10004871", 
    },
    "description" : "won ticket #10004870"
}
// Idempotency is achieved by "fingerprint" value
// Consistency is achieved by optimistic concurrency
```
                                    
```json
# GET  /wallets/:id/transactions/latest
// Returs latest transaction of given wallet
    "id" : "{uuid}"
    "refno" : 2     // a tx sequence number per single wallet
    "amount" : 10.0,             
    "fingerprint" : "TX123A001",
    "labels" : {                 
        "couponId" : "10004871", 
    },
    "description" : "won ticket #10004870",
    "created" : "{timestamp}",
    "oldbalance" : 25.0,
    "newbalance" : 35.0
}
```

```json
# GET  /wallets/:id/balance 
// Returs current balance of given wallet
// Only return a float value; not a json object
``` 

### Run tests
```bash
$ go test ./...  -v
$ go test ./...  -v -tags integration
# if integration tests run; unit tests will not run (using build tags)
# Integration tests use same database with api : walletdb.
# You can use different database for integration tests
$ DB_DATABASE="postgresql://postgres:1234@localhost/testdb" go test ./... -v -tags integration
```

#### Create and Drop database
```bash
go run . db --initdb
go run . db --dropdb
```

#### Environment Variables
```bash
$ DB_DATABASE="postgresql://postgres:1234@localhost/testdb"
$ LOGLEVEL=info
```

### TODO
- In api integrationt tests, only happy path is implemented. Implement the rest.
- Use pgx connection pooling in repository
- Run validations also at wallet service (validations only run at api handlers at the moment) 

### Missing Features
- Get wallet by external id
- Get transaction by fingerprint
- Searching wallets by label
- Searching transactions by label
- Searching transactions by wallet
- Enable api authentication