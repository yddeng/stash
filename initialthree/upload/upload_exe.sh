#!/bin/sh
set -euv

source ./config.sh

# cp config
test -d ${work_dir} || mkdir -p ${work_dir}
cd ${work_dir}
test -d bin && rm -rf bin
mkdir bin

cd bin
# make exe
	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_game/main/node_game.go

  GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/center/main/center.go
	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_dir/main/node_dir.go
	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_login/main/node_login.go
	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_gate/main/node_gate.go
	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_webservice/main/node_webservice.go
	GOOS=${goos_name} GOARCH=${goarch_name} go build -o node_rank ${initialthree_path}/node/node_rank/main/node_rank.go
#
#	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_team/main/node_team.go
#	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_world/main/node_world.go
#	GOOS=${goos_name} GOARCH=${goarch_name} go build ${initialthree_path}/node/node_map/main/node_map.go
#
#	GOOS=${goos_name} GOARCH=${goarch_name} go build -o flykv ${flyfish_path}/server/flykv/main/flykv.go
cd ../

#scp
#if [ "${ssh_used}" = "yes" ];then
#  tar -zcvf bin.tar.gz bin
#  ssh -p ${ssh_port} ${ssh_host} "test -d ${work_dir} || mkdir -p ${work_dir}"
#  scp -r -P ${ssh_port} bin.tar.gz ${ssh_host}:${work_dir}
#  ssh -p ${ssh_port} ${ssh_host} "cd ${work_dir};tar -zxvf bin.tar.gz;rm bin.tar.gz"
#  rm bin.tar.gz
#fi