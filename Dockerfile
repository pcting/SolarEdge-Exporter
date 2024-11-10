FROM golang:1.23 AS build
RUN mkdir /see-build
ADD . /see-build/
WORKDIR /see-build
RUN CGO_ENABLED=0 go build -o solaredge-exporter .

FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /root
COPY --from=build /see-build/solaredge-exporter .
COPY --from=build /see-build/config.yaml /etc/solaredge-exporter/
CMD ["./solaredge-exporter"]
