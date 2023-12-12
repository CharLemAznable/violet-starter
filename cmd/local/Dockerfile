FROM --platform=linux/amd64 busybox:1.36.1

ADD violet-local /app/

WORKDIR /app

VOLUME /violet-configs

EXPOSE 22915 22920

CMD ["/app/violet-local", "-configFile", "/violet-configs/config.toml"]
