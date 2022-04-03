FROM golang:1.18.0-buster as build

RUN apt-get update

COPY . /apitask
RUN cd /apitask && go build -o service main.go

FROM debian:buster-slim

RUN mkdir -p /opt/apitask
COPY --from=build /apitask/docs/swagger.yaml /opt/apitask/docs/swagger.yaml
COPY --from=build /apitask/service /opt/apitask/service
COPY --from=build /apitask/config/config.toml /opt/apitask/config/config.toml

WORKDIR /opt/apitask
CMD ["./service"]