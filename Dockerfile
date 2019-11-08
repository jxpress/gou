FROM golang:latest as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR /opt/app
COPY . /opt/app
RUN go build

# runtime image
FROM alpine
COPY --from=builder /opt/app /opt/app

CMD /opt/app/gou
