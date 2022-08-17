#!/bin/sh

set -ex
cd $(dirname $0)
./copy_template.sh
source inc.sh

tar_name=robot_${goos}.tar

upload_sh="cd ~ && 
( test -d ${remote_workdir} && rm -rf ${remote_workdir}; mkdir -p ${remote_workdir} ) && 
tar -xvf ${tar_name} -C ${remote_workdir} && rm ${tar_name}"

cd .. 
make ${goos}_zip \
	WORK_DIR=${work_dir} \
	GOARCH=${goarch} \
	SERVICE=${service} \
	EXCEL_PATH=${excel_path} \
	QUEST_PATH=${quest_path} 
cd ${work_dir} 
scp -P ${remote_ssh_port} ${tar_name} ${remote_ssh_host}:~ 
ssh -p ${remote_ssh_port} ${remote_ssh_host} "${upload_sh}" 
rm ${tar_name}
