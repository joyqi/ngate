FROM golang:latest as build

COPY ./ /build
WORKDIR /build

RUN CGO_ENABLED=0 go build

FROM alpine:latest

WORKDIR /ngate

COPY --from=build  /build/ngate .
COPY --from=build /build/sample-config.yaml config.yaml

ENTRYPOINT ["./ngate"]