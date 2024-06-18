FROM golang:1.22.4-alpine3.19

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]
