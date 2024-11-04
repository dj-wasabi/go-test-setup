# go-test-setup




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


# 
model/services "function" --> port/in/api.go (Interface) --> adapter/in/http/api used via `cs.uc.function>`


# 
port/out (Interface) "function" --> adapter/out/mongodb "function"


Graceful server shutdown --> https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
