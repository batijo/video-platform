FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /build

COPY --from=builder /src/main .
COPY --from=builder /src/.env .

EXPOSE 8080

CMD ["./main"]