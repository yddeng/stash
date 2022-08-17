#!/bin/sh

set -ex
cd $(dirname $0)
./copy_template.sh
source inc.sh

ssh -p ${remote_ssh_port} ${remote_ssh_host} "cd ${remote_workdir} && ./stop.sh"