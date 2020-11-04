FROM golang:1.13-alpine as build
RUN apk add --update make && rm -rf /var/cache/apk/*
RUN mkdir -vp /workspace
WORKDIR /workspace
ADD . /workspace/
RUN make build

FROM alpine:latest
RUN apk add --update dumb-init && rm -rf /var/cache/apk/*
RUN mkdir -vp /app
WORKDIR /app

COPY --from=build /workspace/bin/fh-geo-svc /app/fh-geo-svc

ENV LISTEN 0.0.0.0:80
EXPOSE 80
ENTRYPOINT [ "dumb-init", "--" ]
CMD [ "/app/fh-geo-svc" ]
