# go-test-setup





# 
model/services "function" --> port/in/api.go (Interface) --> adapter/in/http/api used via `cs.uc.function>`


# 
port/out (Interface) "function" --> adapter/out/mongodb "function"


Graceful server shutdown --> https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
