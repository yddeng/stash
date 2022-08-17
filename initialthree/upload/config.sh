#!/bin/sh
set -euv

initialthree_path="/Users/yidongdeng/go/src/initialthree"  # 需配置 initialthree 路径
config_path="/Users/yidongdeng/svn/TheInitial3_Dev/TheInitial3_Dev_Config" # initialthree 的配置文件目录

work_dir="workspace/initialthree" # 工作目录
ssh_used="no"  # 是否上传

if [ "$1" = "local" ]
then
  goos_name="darwin"
  goarch_name="amd64"
  ip="10.128.2.123"
elif [ "$1" = "212" ]
then
  ssh_used="yes"
  goos_name="linux"
  goarch_name="amd64"
  ip="212.129.131.27"
  ssh_port=59672
  ssh_host="common@212.129.131.27"
elif [ "$1" = "81" ]
then
  config_path="/Users/yidongdeng/svn/TheInitial3_Edition_Config"
  ssh_used="yes"
  goos_name="linux"
  goarch_name="amd64"
  ip="81.69.172.73"
  ssh_port=59672
  ssh_host="common@81.69.172.73"
elif [ "$1" = "150"  ]
then
  config_path="/Users/yidongdeng/svn/TheInitial3_Edition_Config"
  ssh_used="yes"
  goos_name="darwin"
  goarch_name="amd64"
  ip="10.128.2.150"
  ssh_port=22
  ssh_host="feiyu@10.128.2.150"
fi