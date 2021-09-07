# Wallet
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/polarbit/bluelabs-wallet)

## Implementation Notes
- a REST api is provided with minimal endpoints (echo)
- I made use of a sql database (postgres) to achieve consistency
- consistency is achieved by optimistic concurrency and unique keys (via sql database)
- idempotency is achieved by "fingerprint" value (api returns http conflict 409 in case of retry)
- implemented a cli style app (cobra)
- unit tests are written for service package
- integration tests are written for api and db packages

## How to Run
- You need go:1.17, docker and docker-compose
- First run the database server     `$ docker-compose up -d`
- Create the database and tables    `$ go run . db --initdb`
- Run unit tests                    `$ go test ./... -v`
- Run integration tests             `$ go test ./... -v -tags integration`
- Run the api                       `$ go run . api`   
- Repository can also be opened and run on GitPod (free, but login required).

## Enpoints

##### Create Wallet
- `POST /wallets`
- *externalid* value should be unique, enables idempotency
- *labels* is a string dictionary to attach metadata
- response json includes *id* (int)
- request:
```json
{
    "externalId" : "userid-123",
    "labels" : {
        "somekey" : "somevalue"
    }
}
```

##### Get Wallet
- `GET  /wallets/:id`
- response: 
```json
{
    "id": 9180,
    "labels" : {
        "somekey" : "somevalue"
    },
    "externalId" : "userid-123",
    "created" : "2021-09-07T01:42:00"
}
```  
  
##### Add Transaction
- `POST  /wallets/:id/transactions`
- *amount* may be negative or positive (should be <= -1.0 or >=1.0)
- *fingerprint* should be unique, enables idempotency
- *labels* is a string dictionary to attach metadata
- request:
```json
{
    "amount" : 10.0,             
    "fingerprint" : "TX123A001",
    "labels" : { 
        "couponId" : "10004871", 
    },
    "description" : "won ticket #10004870"
}
```
  
##### Get Latest Transaction
- `GET  /wallets/:id/transactions/latest`
- returs latest transaction of the requested wallet
- *refno:* a transaction sequence number per wallet, strengthens consistency
- response: 
```json
{
    "id" : "{uuid}",
    "refno" : 2,
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

##### Get Balance
- `GET  /wallets/:id/balance`
- returs current balance of the requested wallet
- only return a float value; not a json object
```json
35.0
``` 

## Other Notes

#### Run tests
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

### Refactor TODO
- In api integrationt tests, only happy path is implemented. Implement the rest.
- Use pgx connection pooling in repository
- Run validations also at wallet service (validations only run at api handlers at the moment) 

### Missing & Possible Features
- Functional
- Get wallet by external id
- Get transaction by fingerprint
- Search wallets by label
- Search transactions by label
- List transactions by wallet 
- Publish external events (WalletCreated, TransactionCreated, BalanceChanged)
- Technical:
  - enable api authentication
  - enable OpenAPI/Swagger definitions and discovery
  - enable healtcheck and metrics endpoints (may include custom metrics)
  - implement tracing for critical endpoints
