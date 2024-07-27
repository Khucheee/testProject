FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /Customers_kuber ./cmd/main.go
EXPOSE 8080
CMD ["/Customers_kuber"]