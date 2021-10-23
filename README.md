# MongoDB REST API

## Installation

The recommended way to get started using the MongoDB Go driver is by using go modules to install the dependency in
your project. This can be done either by importing packages from `go.mongodb.org/mongo-driver` and having the build
step install the dependency or by explicitly running

```bash
go get go.mongodb.org/mongo-driver/mongo
```

## Usage

Document contains these six lines:

1) _id (ObjectID)- id of the document
2) Company (string) - name of the company
3) Timestamp (string) - pushing time with format: "2006-01-02 15:04:05"
4) Longitude (float64) 
5) Attitude
6) Velocity

### Set Json

To push json in MongoDB, you can write it in test.json with the following fields:

1) "company": "STRING",
2) "longitude": FLOAT64,
3) "attitude": FLOAT64,
4) "velocity": FLOAT64

ID and Timestamp will be generated automatically

### Run

Run the application:

```bash
go run main.go
```

Push json:

```bash
curl -X POST -d @test.json http://localhost:8080/kick
```
