FROM golang:latest as build

COPY ./ /build
WORKDIR /build

RUN CGO_ENABLED=0 go build

FROM alpine:latest

WORKDIR /ngate

COPY --from=build --chmod=0755  /build/ngate .

ENTRYPOINT ["./ngate"]