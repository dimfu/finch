FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .

RUN mkdir /out && \
	CGO_ENABLED=0 GOOS=linux go build -v -o /out/app ./...

FROM scratch
COPY --from=builder /out/app /app

EXPOSE 8080
ENTRYPOINT ["/app"]
