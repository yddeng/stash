#!/bin/sh

set -ex
cd $(dirname $0) 
./copy_template.sh
source inc.sh

./upload_make_config.sh

cd .. && cd ${work_dir}
tar -acvf config.tar config 
scp -P ${remote_ssh_port} config.tar ${remote_ssh_host}:${remote_workdir} 
ssh -p ${remote_ssh_port} ${remote_ssh_host} "cd ${remote_workdir} && rm -rf config && tar -xvf config.tar && rm config.tar"
rm config.tar