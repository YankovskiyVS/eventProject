FROM golang:alpine AS builder

WORKDIR /build/events

ADD go.mod .

COPY . .

RUN go build -o main main.go handle_create.go handle_delete.go handle_update.go postgresql.go http.go typesDB.go message_broker.go

FROM alpine

WORKDIR /build/events

COPY --from=builder /build/events /build/events

EXPOSE 8888

CMD ["/build/events/main"]