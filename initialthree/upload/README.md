# upload

上传工具

### 说明

config.sh            存放所有配置信息
upload_exe.sh        重新编译二进制程序，并上传
upload_source.sh     拷贝资源文件，并上传
upload_install.sh    初始化所有文件
upload_restart.sh    远程开机
upload_stop.sh       远程停机

### 部署

`./upload_install.sh argument` 附带参数，初始化的目标服务器

### 更新

`./upload_exe.sh argument ` 更新程序，附带参数，初始化的目标服务器

`./upload_source.sh argument` 更新资源，附带参数，初始化的目标服务器

### 开关机

`./upload_restart.sh argument` 重启

`./upload_stop.sh argument` 关机