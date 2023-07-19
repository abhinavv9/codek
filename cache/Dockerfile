FROM golang:1.17-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/store

COPY --from=builder /app/cache .

CMD [ "/app/cache", "--port=8080" ]
