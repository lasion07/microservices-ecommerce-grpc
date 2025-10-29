FROM rasadov/ecommerce-base:latest AS build
COPY order order
COPY payment payment
COPY pkg pkg
RUN GO111MODULE=on go build -mod mod -o /go/bin/app ./payment/cmd/payment

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]