#
# BUILD IMAGE
#
FROM --platform=$TARGETPLATFORM golang:1.14.4-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN apk add --update --no-cache git build-base linux-headers

WORKDIR /build

RUN git clone https://github.com/wirepas/c-mesh-api.git && cd c-mesh-api/lib && make

COPY . .

RUN mkdir include lib
RUN cp c-mesh-api/lib/api/* include
RUN cp c-mesh-api/lib/build/mesh_api_lib.a lib/libwirepasmeshapi.a

ENV GO111MODULE=on
ENV TARGET=$TARGETPLATFORM

RUN export PLATFORM=$(echo $TARGET | sed "s/linux\///"); \
    CGO_ENABLED=1 GOOS=linux GOARCH=$PLATFORM \
    go build -a -installsuffix cgo -o wirepas-sink-bridge


#
# RELEASE IMAGE
#
FROM --platform=$TARGETPLATFORM alpine:3.12

WORKDIR /root/
COPY --from=builder /build/wirepas-sink-bridge .

CMD ["./wirepas-sink-bridge"]
