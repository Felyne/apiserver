#!/bin/bash

SERVER="apiserver"
BASE_DIR=$PWD
INTERVAL=2

# 命令行参数，需要手动指定
ARGS=""

function start()
{
	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		echo "$SERVER already running"
		exit 0
	fi

	nohup $BASE_DIR/$SERVER $ARGS &>/dev/null &

	echo "waiting..." &&  sleep $INTERVAL

	# check status
	if [ "`pgrep $SERVER -u $UID`" == "" ];then
		echo "$SERVER start failed"
		exit 1
	fi
}

function status() 
{
	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		echo $SERVER is running
	else
		echo $SERVER is not running
	fi
}

function stop() 
{
	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		kill -15 `pgrep $SERVER -u $UID`
	fi

	echo "waiting..." &&  sleep $INTERVAL

	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		echo "$SERVER stop failed"
		exit 1
	fi
}

function kill_f() 
{
	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		kill -9 `pgrep $SERVER -u $UID`
	fi

	echo "waiting..." &&  sleep $INTERVAL

	if [ "`pgrep $SERVER -u $UID`" != "" ];then
		echo "$SERVER kill failed"
		exit 1
	fi
}

case "$1" in
	'start')
	start
	;;  
	'stop')
	stop
	;;
	'kill')
	kill_f
	;;  
	'status')
	status
	;;  
	'restart')
	stop && start
	;;  
	*)  
	echo "usage: $0 {start|stop|restart|kill|status}"
	exit 1
	;;  
esac
