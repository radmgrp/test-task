FROM golang:1.16 AS builder
ARG goos=linux
ARG goarch=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1000 appuser
WORKDIR /app

COPY --from=builder /app/app .

RUN chown -R appuser /app
USER appuser


EXPOSE 8080

# Command to run the executable
CMD ["./app"]