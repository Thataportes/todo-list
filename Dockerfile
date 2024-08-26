FROM golang:1.23-alpine

WORKDIR /TODO-LIST

COPY . .

RUN go build -o task ./cmd/gotasks/main.go

EXPOSE 8080

CMD ["./task"]
