FROM golang:latest

RUN mkdir -p front

WORKDIR /front

COPY go.mod ./

COPY . .

RUN go build -o main ./cmd

EXPOSE 8082

CMD ["./main"]