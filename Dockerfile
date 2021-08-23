FROM golang:1.17.0 as builder
WORKDIR /go/src/github.com/moosetheory/dumbot/bot
COPY ./ /go/src/github.com/moosetheory/dumbot
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go install -v ./...

FROM alpine:latest
RUN apk --no-cache add tzdata ca-certificates
COPY --from=builder /go/bin/bot /app
CMD ["/app"]
