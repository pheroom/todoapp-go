FROM golang:1.26.5-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/exe .

FROM alpine:3.24
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder --chown=app:app /app/exe /app/exe
USER app
CMD ["/app/exe"]