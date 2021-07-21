# syntax=docker/dockerfile:1
FROM golang:1.16
WORKDIR /app
COPY /app ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app ./
CMD ["./app"]  