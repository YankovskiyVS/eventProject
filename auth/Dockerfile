FROM golang:1.19-alpine

WORKDIR /github/eventProject/auth

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]