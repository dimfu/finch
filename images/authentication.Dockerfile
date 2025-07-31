FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY authentication /build/authentication

WORKDIR /build/authentication

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /bin/authentication .

FROM scratch
WORKDIR /app
COPY --from=builder /bin/authentication .

EXPOSE 8080
ENTRYPOINT ["/app/authentication"]
