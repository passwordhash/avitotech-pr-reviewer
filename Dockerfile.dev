FROM golang:1.24-alpine3.20 AS builder

RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
 && apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /server.app cmd/http/main.go

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=appuser:appgroup --chmod=755 /server.app /server.app
COPY --from=builder --chown=appuser:appgroup /app/configs /configs

USER appuser

EXPOSE 8080

CMD ["/server.app"]
