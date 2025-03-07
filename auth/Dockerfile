FROM golang:alpine AS builder

WORKDIR /build/auth

ADD go.mod .

COPY . .

RUN go build -o main main.go handlers.go user.go

FROM alpine

WORKDIR /build/auth

COPY --from=builder /build/auth /build/auth

EXPOSE 8081

CMD ["/build/auth/main"]