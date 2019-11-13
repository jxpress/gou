FROM golang:latest as builder

RUN apt-get update && apt-get install build-essential -y

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR /opt/app
COPY . /opt/app
RUN go build -a -ldflags '-linkmode external -extldflags "-static"'

# runtime image
FROM alpine
COPY --from=builder /opt/app/ /opt/app
COPY create_table.sql /opt/app/create_table.sql

CMD /opt/app/gou
