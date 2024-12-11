FROM golang:1.23 AS build
WORKDIR /go/src
ADD . .

ENV CGO_ENABLED=0
RUN go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
RUN go generate /go/src/internal/adapter/in/http/api/server-generator.go
RUN go generate /go/src/internal/core/domain/model/model-generator.go
RUN go generate /go/src/internal/core/port/in/in-generator.go
RUN go build -o /go/src/api ./cmd/api 

FROM scratch AS runtime

ENV GIN_MODE=release

COPY --from=build /go/src/api ./
COPY --from=build /go/src/config.yaml ./

EXPOSE 8888/tcp
ENTRYPOINT ["./api"]