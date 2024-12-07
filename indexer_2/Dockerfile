FROM golang:1.21.12-alpine

RUN mkdir app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o indexer .

CMD ["./indexer", "/app/enron_mail_20110402/maildir/"]