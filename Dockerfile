FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/website /opt/bin/website

CMD ["/opt/bin/website"]