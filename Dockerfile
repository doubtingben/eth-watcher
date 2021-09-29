FROM golang:1.16 as builder

ADD ./ /app/
WORKDIR /app/
RUN go build

FROM google/cloud-sdk
COPY --from=builder /app/eth-watcher /usr/local/bin/eth-watcher
