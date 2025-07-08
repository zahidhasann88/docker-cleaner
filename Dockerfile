FROM alpine:3.20 as runtime

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-cleaner .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /root/

COPY --from=builder /app/docker-cleaner .

RUN chown appuser:appuser docker-cleaner

USER appuser

ENTRYPOINT ["./docker-cleaner"]