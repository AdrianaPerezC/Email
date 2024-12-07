FROM golang:1.21.12-alpine

RUN mkdir app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

ARG PORT=:8080 \
    GIN_MODE=release

ENV PORT=${PORT} \
    GIN_MODE=${GIN_MODE}

COPY . .

RUN go build -o main .

CMD ["./main"]