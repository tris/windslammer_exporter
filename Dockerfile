FROM golang:1.21-alpine AS builder
LABEL org.opencontainers.image.authors="Tristan Horn <tristan+docker@ethereal.net>"
WORKDIR /app
RUN apk add --no-cache upx
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o windslammer_exporter .
RUN upx --lzma windslammer_exporter

FROM scratch
COPY --from=builder /app/windslammer_exporter /windslammer_exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/windslammer_exporter"]