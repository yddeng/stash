#!/bin/sh
set -euv

flyfishPath="/Users/yidongdeng/go/src/github.com/sniperHW/flyfish"
goos_name="linux"
goarch_name="amd64"
ip="212.129.131.27"
ssh_port=59672
ssh_host="common@212.129.131.27"

test -d ../workspace/flyfish || mkdir ../workspace/flyfish

cd ../workspace/flyfish
## make exe
GOOS=${goos_name} GOARCH=${goarch_name} go build ${flyfishPath}/server/app/flygate.go
GOOS=${goos_name} GOARCH=${goarch_name} go build ${flyfishPath}/server/app/flykv.go
GOOS=${goos_name} GOARCH=${goarch_name} go build ${flyfishPath}/server/app/flypd.go

test -f flygate_config.toml || cp ${flyfishPath}/server/app/flygate_config.toml.template ./flygate_config.toml
test -f flykv_config.toml || cp ${flyfishPath}/server/app/flykv_config.toml.template ./flykv_config.toml
test -f flypd_config.toml || cp ${flyfishPath}/server/app/flypd_config.toml.template ./flypd_config.toml
test -f meta.json || cp ../../flyfish/meta.json ./meta.json
test -f deployment.json || cp ../../flyfish/deployment.json ./deployment.json


cp ../../flyfish/start.sh ./start.sh
cp ../../flyfish/stop.sh ./stop.sh
chmod +x start.sh
chmod +x stop.sh


#cd ../
#tar -zcvf flyfish.tar.gz flyfish
#ssh -p ${ssh_port} ${ssh_host} "test -d workspace || mkdir -p workspace"
#scp -r -P ${ssh_port} flyfish.tar.gz ${ssh_host}:workspace
#ssh -p ${ssh_port} ${ssh_host} "cd workspace;tar -zxvf flyfish.tar.gz;rm flyfish.tar.gz"
#rm flyfish.tar.gz