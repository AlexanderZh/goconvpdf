FROM golang:1.18 AS build
RUN apt install -y ca-certificates
WORKDIR /go/src

ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GOARCH=amd64

COPY ./ /go/src/
RUN go mod tidy



RUN go build
RUN ls

FROM alpine:latest AS runtime

RUN set -xe \
    && apk add --no-cache --purge -uU \
        curl icu-libs unzip zlib-dev musl \
        mesa-gl mesa-dri-swrast \
        libreoffice libreoffice-base libreoffice-lang-uk \
        ttf-freefont ttf-opensans ttf-inconsolata \
	ttf-liberation ttf-dejavu \
        libstdc++ dbus-x11 \
    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
    && apk add --no-cache -U \
	ttf-font-awesome ttf-mononoki ttf-hack \
    && rm -rf /var/cache/apk/* /tmp/*

COPY --from=build /go/src/goconvpdf ./
RUN ls
EXPOSE 8080/tcp
ENTRYPOINT ["./goconvpdf"]
