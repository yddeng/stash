#!/bin/sh

set -ex
cd $(dirname $0)
./copy_template.sh
source inc.sh

upload_bin_sh="cd ${remote_workdir} && tar -xvf robot_bin.tar && rm robot_bin.tar"

cd ..
make ${goos} WORK_DIR=${work_dir} GOARCH=${goarch} 
cd ${work_dir} 
tar -acvf robot_bin.tar ./bin/robot 
scp -P ${remote_ssh_port} robot_bin.tar ${remote_ssh_host}:${remote_workdir} 
ssh -p ${remote_ssh_port} ${remote_ssh_host} "${upload_bin_sh}"
rm robot_bin.tar