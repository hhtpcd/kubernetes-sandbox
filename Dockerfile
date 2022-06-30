FROM golang:1.18-bullseye as builder

WORKDIR /app

ADD . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o server

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/server /server

ENTRYPOINT [ "/server" ]
