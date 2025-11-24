FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 go build -v -o main .
RUN ls -l main

EXPOSE 8080

CMD ["./main"]

