#!/bin/sh
set -euv

source ./config.sh

# cp config
test -d ${work_dir} || mkdir -p ${work_dir}
cd ${work_dir}
test -d configs && rm -rf ./configs
mkdir configs

cp ${initialthree_path}/sql.sql ./
cp -r ${config_path}/* ./configs/

#scp
if [ "${ssh_used}" = "yes" ];then
  tar -zcvf configs.tar.gz configs sql.sql
  ssh -p ${ssh_port} ${ssh_host} "test -d ${work_dir} || mkdir -p ${work_dir}"
  scp -r -P ${ssh_port} configs.tar.gz ${ssh_host}:${work_dir}
  ssh -p ${ssh_port} ${ssh_host} "cd ${work_dir};rm -rf configs;tar -zxvf configs.tar.gz;rm configs.tar.gz"
  rm configs.tar.gz
fi

