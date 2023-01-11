
FROM golang:alpine AS builder

RUN mkdir /app

ADD . /app

WORKDIR /app/cmd

RUN go build -o main .

FROM alpine

COPY --from=builder /app/cmd/main /app/

CMD ["/app/main"]
