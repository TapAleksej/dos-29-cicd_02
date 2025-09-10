FROM golang:1.21-alpine
WORKDIR /app
RUN go mod init ecommerce
#RUN go mod tidy
RUN go get github.com/lib/pq
COPY . .
RUN go build -o app main.go
ENTRYPOINT ["./app"]

