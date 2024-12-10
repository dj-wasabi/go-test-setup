# go-test-setup

A personal project to learn golang and try to set it up in a hexagon architecture setup. Endpoints and the models/objects are generated via an OpenAPI spec file. It is a very basic API that has Authentication, creating or users and organisation possibilities. This isn't production ready code btw.

Some (Prometheus) metrics are collected and are displayed with the `/v1/metrics` endpoint.




# Migrations

```sh
$ migrate create -seq -ext=.json -dir=./migrations <my_reason_to_create_a_migration>
```

Executing to create:
```sh
$ migrate -path=./migrations -database="mongodb://user@pass:localhost:27017/mongo-golang-test" up
```

Or to rollback:
```sh
$ migrate -path=./migrations -database="mongodb://user@pass:localhost:27017/mongo-golang-test" down
```

