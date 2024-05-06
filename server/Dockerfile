FROM golang:latest

RUN mkdir -p back

WORKDIR /back

COPY go.mod ./

COPY . .

RUN go build -o main ./cmd

EXPOSE 8083

CMD ["./main"]