#!/bin/bash

BASEDIR="/usr/local/stathub"
PIDFILE="log/stathub.pid"

cd $BASEDIR

start() {
    echo "starting"
    if [ ! -f $PIDFILE ]; then
        ./stathub -c conf/stathub.conf
        echo "ok"
        echo "----------------------------------------------------"
        echo "| Please open the below URL on your PC browser:    |"
        echo "| https://this-server-ip:15944                     |"
        echo "| Default password: likexian                       |"
        echo "|                                                  |"
        echo "| Feedback: https://github.com/likexian/stathub-go |"
        echo "| Thank you for your using, By Li Kexian           |"
        echo "| StatHub, Apache License, Version 2.0             |"
        echo "----------------------------------------------------"
    fi
}

stop() {
    echo "stopping"
    if [ -f $PIDFILE ]; then
        kill -9 `cat $PIDFILE`
        rm -rf $PIDFILE
        echo "ok"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
esac
