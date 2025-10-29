FROM rasadov/ecommerce-base:latest AS build
COPY product product
COPY pkg pkg
RUN GO111MODULE=on go build -mod mod -o /go/bin/app ./product/cmd/product

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]