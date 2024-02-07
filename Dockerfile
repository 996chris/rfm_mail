# Stage 1: Build
FROM golang:1.21.5-alpine3.18 as builder
WORKDIR /app
COPY . ./source
RUN go build -C source -v -o /service

# Stage 2: Final Image
FROM alpine:3.18
COPY --from=builder /service /service
ENTRYPOINT [ "/service" ]

