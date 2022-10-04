#
# Build go project
#
FROM golang:1.19.1-alpine as go-builder

WORKDIR /go/src/test/gohttp

COPY ./gohttp/ .

#RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gohttp *.go

#
# Runtime container
#
FROM alpine:latest  

RUN mkdir -p /app && \
    addgroup -S app && adduser -S app -G app && \
    chown app:app /app

WORKDIR /app

COPY --from=go-builder /go/src/test/gohttp .

USER app

CMD ["./gohttp"]  
