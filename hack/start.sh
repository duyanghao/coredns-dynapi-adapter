#!/bin/bash

server="./coredns-dynapi-adapter"
let item=0
item=`ps -ef | grep $server | grep -v grep | wc -l`

if [ $item -eq 1 ]; then
	echo "The coredns-dynapi-adapter is running, shut it down..."
	pid=`ps -ef | grep $server | grep -v grep | awk '{print $2}'`
	kill -9 $pid
fi

echo "Start coredns-dynapi-adapter now ..."
make
./build/pkg/cmd/coredns-dynapi-adapter/coredns-dynapi-adapter  >> coredns-dynapi-adapter.log 2>&1 &
