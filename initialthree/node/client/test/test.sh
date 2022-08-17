#!/bin/bash

for((i=1;i<=100;i++));
do
#nohup ./createRole 10.128.2.123:9201 $(expr $i + 1000) &
nohup ./enterMap 10.128.2.123:9201 $(expr $i + 1000) &
done