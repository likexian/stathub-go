#!/bin/bash

VERSION="0.102.4"
STATHUB_URL="https://github.com/likexian/stathub-go/releases/download/v${VERSION}/stathub.$(uname -m).tar.gz"

BASEDIR="/usr/local/stathub"

[ $(id -u) -ne 0 ] && sudo="sudo" || sudo=""
id -u nobody >/dev/null 2>&1
if [ $? -ne 0 ]; then
    $sudo groupadd nogroup
    useradd -g nogroup nobody -s /bin/false
fi

$sudo mkdir -p $BASEDIR
$sudo chown -R $(id -u -n):$(id -g -n) $BASEDIR
if [ ! -d $BASEDIR ]; then
    echo "Unable to create dir $BASEDIR and chown to current user, Please manual do it"
    exit 1
fi
cd $BASEDIR

command_exists() {
    type "$1" &> /dev/null
}

if command_exists wget; then
    wget --no-check-certificate $STATHUB_URL
elif command_exists curl; then
    curl --insecure -O $STATHUB_URL
else
    echo "Unable to find curl or wget, Please install and try again."
    exit 1
fi

if [ ! -f stathub.$(uname -m).tar.gz ]; then
    echo "Unable to get server file, Please manual download it"
    exit 1
fi

tar zxf stathub.$(uname -m).tar.gz
chmod +x stathub service
[ ! -d conf ] && $sudo mkdir $BASEDIR/conf
if [ ! -f conf/stathub.conf ]; then
    $sudo ./stathub -c conf/stathub.conf --init-server
fi
$sudo mkdir $BASEDIR/pkgs
mv stathub.*.tar.gz $BASEDIR/pkgs

if [ -z "$(grep stathub /etc/rc.local)" ]; then
    $sudo sh -c "echo \"cd $BASEDIR; rm -f log/stathub.pid; ./service start >>$BASEDIR/log/stathub.log 2>&1\" >> /etc/rc.local"
fi

echo "----------------------------------------------------"
echo "| Server install successful, Please start it using |"
echo "| ./service {start|stop|restart}                   |"
echo "| Now it will automatic start                      |"
echo "|                                                  |"
echo "| Feedback: https://github.com/likexian/stathub-go |"
echo "| Thank you for your using, By Li Kexian           |"
echo "| StatHub, Apache License, Version 2.0             |"
echo "----------------------------------------------------"

$sudo ./service start
