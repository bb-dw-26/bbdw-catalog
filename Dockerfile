FROM golang:1.27.3-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /catalog ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache ca-certificates && adduser -D app
COPY --from=build /catalog /usr/local/bin/catalog
USER app
EXPOSE 3002
CMD ["catalog"]
