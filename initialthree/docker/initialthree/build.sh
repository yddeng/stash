#!/bin/sh
set -euv

initialthree_path="/Users/yidongdeng/go/src/initialthree"  # 需配置 initialthree 路径
config_path="/Users/yidongdeng/svn/TheInitial3_Dev/TheInitial3_Dev_Config" # initialthree 的配置文件目录
goos_name="linux"
goarch_name="amd64"


test -d bin && rm -rf bin
mkdir bin
cd bin
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/node/node_game/main/node_game.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/center/main/center.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/node/node_dir/main/node_dir.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/node/node_login/main/node_login.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/node/node_gate/main/node_gate.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${initialthree_path}/node/node_webservice/main/node_webservice.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" -o node_rank ${initialthree_path}/node/node_rank/main/node_rank.go
cd ..


test -d configs && rm -rf ./configs
mkdir configs
cp -r ${config_path}/* ./configs/

test -f config.toml ||  cp ../../upload/template/config.toml.template config.toml