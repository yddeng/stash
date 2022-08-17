#!/bin/sh

set -ex
cd $(dirname $0) 
./copy_template.sh
source inc.sh

cd ..
make copy_config copy_shell \
    WORK_DIR=${work_dir} \
    SERVICE=${service} \
    EXCEL_PATH=${excel_path} \
    QUEST_PATH=${quest_path}