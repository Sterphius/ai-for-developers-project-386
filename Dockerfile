FROM golang:1.25-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY . .

RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/api /usr/local/bin/api
COPY --from=build /src/openapi /app/openapi

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]
