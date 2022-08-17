#Dockerfile指令说明简洁版：
#FROM
#构建镜像基于哪个镜像
#MAINTAINER
#镜像维护者姓名或邮箱地址
#RUN
#构建镜像时运行的指令
#CMD
#运行容器时执行的shell环境
#VOLUME
#指定容器挂载点到宿主机自动生成的目录或其他容器
#USER
#为RUN、CMD、和 ENTRYPOINT 执行命令指定运行用户
#WORKDIR
#为 RUN、CMD、ENTRYPOINT、COPY 和 ADD 设置工作目录，就是切换目录
#HEALTHCHECH
#健康检查
#ARG
#构建时指定的一些参数
#EXPOSE
#声明容器的服务端口（仅仅是声明）
#ENV
#设置容器环境变量
#ADD
#拷贝文件或目录到容器中，如果是URL或压缩包便会自动下载或自动解压
#COPY
#拷贝文件或目录到容器中，跟ADD类似，但不具备自动下载或解压的功能
#ENTRYPOINT
#运行容器时执行的shell命令


# 多阶段构建

# 基础镜像 step2
# 编译项目
FROM golang:1.17.7 AS builder
# 设置工作区
# 由于项目用的govender方式，这里需要指定GOPATH路径
WORKDIR /go/src/initialthree
# 拷贝项目文件到工作目录
COPY . .
# 环境变量
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=off
# 编译可执行的二进制文件
# -ldflags：用于传递每个go工具链接调用的参数。
#   -s: 省略符号表和调试信息
#   -w: 省略DWARF符号表
#RUN  go build -ldflags="-s -w" -o center center/main/center.go
#RUN  go build -ldflags="-s -w" -o node_dir node/node_dir/main/node_dir.go
#RUN  go build -ldflags="-s -w" -o node_login node/node_login/main/node_login.go
#RUN  go build -ldflags="-s -w" -o node_gate node/node_gate/main/node_gate.go
#RUN  go build -ldflags="-s -w" -o node_game node/node_game/main/node_game.go
RUN  go build -ldflags="-s -w" -o node_webservice node/node_webservice/main/node_webservice.go
#RUN  go build -ldflags="-s -w" -o node_rank node/node_rank/main/node_rank.go
#RUN  go build -ldflags="-s -w" -o node_conflictzone node/node_conflictzone/main/main.go




# 程序镜像附带taegets step
# 当我们在构建镜像的时候就可以通过 --target 参数来明确指定要构建的 Targets
# --tag 为镜像指定一个标签


# 编译镜像 docker build --target=center --tag=center:1.0 .
# 启动镜像 docker run -d -p 40500:40500 center:1.0
# 进入容器 docker exec -it 容器ID sh (对于基于alpine的容器来说,没有 /bin/bash)
FROM alpine AS center
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/center .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./center"]
# 提供默认参数
CMD ["localhost:40500"]



# node_dir
FROM alpine AS node_dir
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_dir .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_dir"]
# 提供默认参数
CMD ["config.toml", "Dir", "0"]



# node_login
FROM alpine AS node_login
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_login .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_login"]
# 提供默认参数
CMD ["config.toml", "Login", "0"]



# node_gate
FROM alpine AS node_gate
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_gate .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_gate"]
# 提供默认参数
CMD ["config.toml", "Gate", "0"]



# node_game
FROM alpine AS node_game
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_game .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_game"]
# 提供默认参数
CMD ["config.toml", "Game", "0"]



# node_webservice
FROM alpine AS node_webservice
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_webservice .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_webservice"]
# 提供默认参数
CMD ["config.toml", "WebService", "0"]



# node_rank
FROM alpine AS node_rank
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_rank .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_rank"]
# 提供默认参数
CMD ["config.toml", "Rank", "0"]



# node_conflictzone
FROM alpine AS node_conflictzone
WORKDIR /
# 复制在builder阶段的时区，日志就显示的是北京时间了
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 从基础镜像中拷贝执行程序
COPY --from=builder /go/src/initialthree/node_conflictzone .
# 拷贝配置文件
COPY --from=builder /go/src/initialthree/upload/template/config.toml.template config.toml
# 容器进入时执行的命令
ENTRYPOINT ["./node_conflictzone"]
# 提供默认参数
CMD ["config.toml", "ConflictZone", "0"]