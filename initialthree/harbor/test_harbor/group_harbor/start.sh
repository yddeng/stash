#!/bin/sh
nohup go run ../../../center/center.go localhost:8012 initialthree > center.log 2>&1 &