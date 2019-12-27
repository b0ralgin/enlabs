The main goal of this test task is a develop the application for processing the incoming requests from the 3d-party providers.
The application must have an HTTP URL to receive incoming POST requests.
To receive the incoming POST requests the application must have an HTTP URL endpoint.
 
#Initial setup

##Building application
```shell script
go build ./bin/enlabs ./cmd/main.go
```
##Set the environment variables**
```shell script
#binding of http server
ENLABS_ADDR=0.0.0.0:8080
#level of logger messages
ENLABS_LOG_LEVEL=debug
#DB connection string 
ENLABS_DSN=host=localhost user=postgres dbname=postgres password=postgres sslmode=disable
#period of balance corredction 
ENLABS_PERIOD=5
```

##Make migration 
```shell script
./bin/enlabs migrate 
```
#Solution
##Processing and saving incoming requests.

Start http server: 
```shell script
./bin/enlabs server 
```
There are 2 endpoints: 

1. /transaction - enpoint for adding new transaction
```http request
POST /transaction HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Source-Type: server

{"state": "win", "amount": "10.15", "transactionId": "sdf32fsdfdsfw3f3"}
```
2. /balance - shows current account balance
```http request
GET /balance HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

##Post-processing
Start corrector service: 
```shell script
./bin/enlabs corrector
```
Each N minutes (which set as ENLABS_PERIOD in the environment variables) latest odd records will be canceled and balance will be corrected.

#Docker 
There are a dockerfile to build docker container and docker-compose file to run the solution at one step
##build container
```shell script
docker build -t enlabs .
```

##run docker-compose
```shell script
docker-compose up -d db
docker-compose up migrate
docker-compose up -d server corrector
```