#!/bin/sh
set -euv
source ./config.sh

ssh -p ${ssh_port} ${ssh_host} "cd ${work_dir}; ./stop.sh; "