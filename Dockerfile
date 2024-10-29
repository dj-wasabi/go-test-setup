FROM golang:1.23 AS build
WORKDIR /go/src
ADD . .

ENV CGO_ENABLED=0
RUN go build -o /go/src/api ./cmd/api 

FROM scratch AS runtime

ENV GIN_MODE=release

COPY --from=build /go/src/api ./
COPY --from=build /go/src/config.yaml ./

EXPOSE 8888/tcp
ENTRYPOINT ["./api"]