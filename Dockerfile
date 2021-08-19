FROM golang:1.16-alpine AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

FROM alpine
RUN apk add poppler-utils
COPY --from=builder /main ./
COPY --from=builder /app/config.yml ./
ENTRYPOINT ["./main"]
EXPOSE 8080