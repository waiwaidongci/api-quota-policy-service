FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/api-quota-service ./cmd/api-quota-service
FROM alpine:3.20
COPY --from=build /out/api-quota-service /usr/local/bin/api-quota-service
EXPOSE 8084
ENTRYPOINT ["/usr/local/bin/api-quota-service"]
