FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o sponsors-webhook .

FROM alpine:3.24
WORKDIR /app

RUN adduser -D -H -u 10001 app
COPY ./static /app/static/
COPY --from=builder /app/sponsors-webhook .
USER app

EXPOSE 1928
CMD ["./sponsors-webhook"]
