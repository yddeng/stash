# 当我们在构建镜像的时候就可以通过 --target 参数来明确指定要构建的 Targets
# --tag 为镜像指定一个标签


# 编译镜像 docker build --target=center --tag=center:1.0 .
# 启动镜像 docker run -d -p 40500:40500 center:1.0
# 进入容器 docker exec -it 容器ID sh (对于基于alpine的容器来说,没有 /bin/bash)
FROM alpine AS flygate

WORKDIR /flyfish

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 从基础镜像中拷贝执行程序
COPY ./flygate .
# 拷贝配置文件
COPY ./flygate_config.toml .
# 容器进入时执行的命令
ENTRYPOINT ["./flygate"]
# 提供默认参数
CMD ["-config=./flygate_config.toml"]



FROM alpine AS flypd

WORKDIR /flyfish

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 从基础镜像中拷贝执行程序
COPY ./flypd .
# 拷贝配置文件
COPY ./flypd_config.toml .
COPY ./meta.json .
COPY ./deployment.json .
# 容器进入时执行的命令
ENTRYPOINT ["./flypd"]
# 提供默认参数
CMD ["-config=./flypd_config.toml","-raftcluster=1@1@http://localhost:8811@127.0.0.1:10050@voter"]



FROM alpine AS flykv

WORKDIR /flyfish

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 从基础镜像中拷贝执行程序
COPY ./flykv .
# 拷贝配置文件
COPY ./flykv_config.toml .
# 容器进入时执行的命令
ENTRYPOINT ["./flykv"]
# 提供默认参数
CMD ["-config=./flykv_config.toml","-id=1"]