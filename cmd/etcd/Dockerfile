FROM --platform=linux/amd64 busybox:1.36.1

ADD violet-etcd /app/

WORKDIR /app

VOLUME /violet-configs

EXPOSE 22915 22920

CMD ["/app/violet-etcd", "-configFile", "/violet-configs/config.toml"]
