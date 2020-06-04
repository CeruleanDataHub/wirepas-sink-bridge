#
# BUILD IMAGE
#
FROM golang:1.14.4-alpine AS builder

RUN apk add --update --no-cache git build-base linux-headers

WORKDIR /build

RUN git clone https://github.com/wirepas/c-mesh-api.git && cd c-mesh-api/lib && make

COPY . .

ENV GO111MODULE=on

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o wirepas-sink-bridge


#
# RELEASE IMAGE
#
FROM alpine:3.12

WORKDIR /root/
COPY --from=builder /build/wirepas-sink-bridge .

CMD ["./wirepas-sink-bridge"]