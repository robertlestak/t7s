FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/t7s cmd/t7s/*.go

FROM alpine:latest as app

COPY --from=builder /bin/t7s /bin/t7s

RUN chmod +x /bin/t7s

WORKDIR /data
ENTRYPOINT ["/bin/t7s"]