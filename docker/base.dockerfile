FROM golang:1.23-alpine3.20
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/rasadov/EcommerceAPI
COPY go.mod go.sum ./
RUN go mod download