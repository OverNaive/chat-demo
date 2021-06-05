FROM golang:1.16-buster as build

ENV GOPROXY="https://goproxy.io"

WORKDIR /home/app

COPY ./src .

# build
RUN go mod download \
&& go build -o chat *.go


FROM debian:buster-slim as prod

WORKDIR /home/app

# binary file
COPY --from=build /home/app/chat /home/app/chat

# supervisor setting
COPY ./supervisor.conf /etc/supervisor/conf.d/supervisor.conf

RUN chmod +x chat \
&& apt-get update \
# install tools
&& apt-get install -y --no-install-recommends supervisor \
# clean up for smaller size
&& apt-get autoclean \
&& rm -rf /var/lib/apt/lists/

EXPOSE 8888

CMD [ "supervisord", "-c", "/etc/supervisor/supervisord.conf" ]