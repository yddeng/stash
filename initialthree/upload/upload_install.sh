#!/bin/sh
set -euv

source ./config.sh

./upload_source.sh $1
./upload_exe.sh $1

# template
test -d ${work_dir} || mkdir -p ${work_dir}
test -f ${work_dir}/config.toml ||  cp template/config.toml.template ${work_dir}/config.toml
test -f ${work_dir}/start.sh ||  cp template/start.sh.template ${work_dir}/start.sh
test -f ${work_dir}/stop.sh ||  cp template/stop.sh.template ${work_dir}/stop.sh

# 替换对外IP
go run tool/replace_external.go ${work_dir}/config.toml ${ip}

# chmod
cd ${work_dir}
chmod +x start.sh
chmod +x stop.sh

#scp
if [ "${ssh_used}" = "yes" ];then
  tar -zcvf template.tar.gz config.toml start.sh stop.sh
  ssh -p ${ssh_port} ${ssh_host} "test -d ${work_dir} || mkdir -p ${work_dir}"
  scp -r -P ${ssh_port} template.tar.gz ${ssh_host}:${work_dir}
  ssh -p ${ssh_port} ${ssh_host} "cd ${work_dir};tar -zxvf template.tar.gz;rm template.tar.gz"
  rm template.tar.gz
fi
