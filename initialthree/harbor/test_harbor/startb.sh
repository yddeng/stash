#!/bin/zsh
nohup go run ../../center/main/center.go localhost:8011 initialthree > center3.log 2>&1 &
nohup go run ../../harbor/harbor.go localhost:8011@localhost:8012 2.255.1 localhost:9102 initialthree > harbor2.log 2>&1 &