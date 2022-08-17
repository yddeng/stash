#!/bin/sh
set -euv

flyfishPath="/Users/yidongdeng/go/src/github.com/sniperHW/flyfish"
goos_name="linux"
goarch_name="amd64"

## make exe
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${flyfishPath}/server/app/flygate.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${flyfishPath}/server/app/flykv.go
GOOS=${goos_name} GOARCH=${goarch_name} go build -ldflags="-s -w" ${flyfishPath}/server/app/flypd.go

test -f flygate_config.toml || cp ${flyfishPath}/server/app/flygate_config.toml.template ./flygate_config.toml
test -f flykv_config.toml || cp ${flyfishPath}/server/app/flykv_config.toml.template ./flykv_config.toml
test -f flypd_config.toml || cp ${flyfishPath}/server/app/flypd_config.toml.template ./flypd_config.toml
test -f meta.json || cp ../../upload/flyfish/meta.json ./meta.json
test -f deployment.json || cp ../../upload/flyfish/deployment.json ./deployment.json