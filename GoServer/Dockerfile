FROM golang:1.24.1

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
