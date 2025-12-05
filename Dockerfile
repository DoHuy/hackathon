FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hackathon-app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/hackathon-app .
COPY env.ini .
COPY docs/ ./docs/
RUN mkdir -p /tmp
EXPOSE 8080
CMD ["./hackathon-app"]