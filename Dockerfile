FROM golang:alpine AS build

WORKDIR /build
COPY . .

RUN go build -v ./cmd/jnsd

FROM alpine

COPY --from=build /build/jnsd /usr/local/bin/jnsd

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/jnsd"]
