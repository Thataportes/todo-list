FROM golang:1.23-alpine

WORKDIR /TODO-LIST

COPY . .

RUN go build -o cmd/gotasks/main ./cmd/gotasks

EXPOSE 8080

CMD ["./cmd/gotasks/main"]
