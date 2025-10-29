FROM rasadov/ecommerce-base:latest AS build
COPY account account
COPY product product
COPY order order
COPY recommender recommender
COPY payment payment
COPY graphql graphql
COPY pkg pkg
RUN GO111MODULE=on go build -mod mod -o /go/bin/app ./graphql/cmd/graphql

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]