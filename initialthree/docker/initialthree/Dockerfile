
FROM alpine AS center

WORKDIR /initialthree

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 从基础镜像中拷贝执行程序
COPY ./bin/center .
# 拷贝配置文件
COPY ./config.toml .
# 容器进入时执行的命令
ENTRYPOINT ["./center"]
# 提供默认参数
CMD [":40500"]
EXPOSE 40500


# node_dir
FROM alpine AS node_dir
WORKDIR /initialthree
# 复制在builder阶段的时区，日志就显示的是北京时间了
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
# 从基础镜像中拷贝执行程序
COPY ./bin/node_dir .
# 拷贝配置文件
COPY ./config.toml .
# 容器进入时执行的命令
ENTRYPOINT ["./node_dir"]
# 提供默认参数
CMD ["config.toml", "Dir", "0"]
EXPOSE 40501
EXPOSE 40101
