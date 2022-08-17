#!/bin/sh
set -euv
source ./config.sh

ssh -p ${ssh_port} ${ssh_host} "cd ${work_dir}; ./stop.sh;sleep 1s;./start.sh;ps -eo pid,lstart,etime,cmd | grep initialthree "