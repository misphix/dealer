# build stage
FROM golang:1.19 AS builder
WORKDIR /go/src/dealer
COPY . .
RUN make

# final stage
FROM ubuntu:22.04
WORKDIR /root/
COPY --from=builder /go/src/dealer/dealer .
COPY --from=builder /go/src/dealer/config/config.yaml ./config/config.yaml
ENTRYPOINT ["./dealer"]