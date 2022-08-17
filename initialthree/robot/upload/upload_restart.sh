#!/bin/sh

set -ex
cd $(dirname $0)
./copy_template.sh
source inc.sh

ssh -p ${remote_ssh_port} ${remote_ssh_host} "cd ${remote_workdir} && ulimit -n 10000 && ./stop.sh && ./start.sh"