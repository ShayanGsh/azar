
FROM golang:alpine AS builder

RUN mkdir /app

ADD . /app

WORKDIR /app/cmd

RUN go build -o main .

FROM alpine

COPY --from=builder /app/cmd/main /app/
COPY --from=builder /app/azar-db/migrations /app/migrations

RUN chmod +x /app/main

WORKDIR /app

CMD ["./main"]