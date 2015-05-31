#!/bin/sh

VERSION="0.10.2"
STATHUB_URL="https://github.com/likexian/stathub-go/releases/download/v${VERSION}/server_$(uname -m).tar.gz"

[ $UID -ne 0 ] && sudo="sudo" || sudo=""
$sudo mkdir -p /var/stathub
$sudo chown -R $USER:$USER /var/stathub
if [ ! -d /var/stathub ]; then
    echo "Unable to create dir /var/stathub and chown to current user, Please manual do it"
    exit 1
fi
cd /var/stathub

if [ -f $(which wget) ]; then
    wget --no-check-certificate $STATHUB_URL
elif [ -f $(which curl) ]; then
    curl --insecure -O $STATHUB_URL
else
    echo "Unable to find curl or wget, Please install and try again."
    exit 1
fi

if [ ! -f server_$(uname -m).tar.gz ]; then
    echo "Unable to get server file, Please manual download it"
    exit 1
fi
tar zxf server_$(uname -m).tar.gz
rm -rf server_$(uname -m).tar.gz

if [ -z "$(grep stathub /etc/rc.local)" ]; then
    echo "cd /var/stathub; ./service start >/var/stathub/data/stathub.log 2>&1" >> /etc/rc.local
fi

echo "----------------------------------------------------"
echo "| Server install sucessful, Please start it using  |"
echo "| ./service {start|stop|restart}                   |"
echo "| Now it will automatic start                      |"
echo "|                                                  |"
echo "| Feedback: https://github.com/likexian/stathub-go |"
echo "| Thank you for your using, By Li Kexian           |"
echo "| StatHub, Apache License, Version 2.0             |"
echo "----------------------------------------------------"

./service start

KEY=$(grep key server.json | cut -d \" -f 4)
STATHUB_CLIENT_URL=https://127.0.0.1:15944/node?key=$KEY
if [ -f $(which wget) ]; then
    wget --no-check-certificate -O - $STATHUB_CLIENT_URL | sh
elif [ -f $(which curl) ]; then
    curl --insecure $STATHUB_CLIENT_URL | sh
else
    echo "Unable to find curl or wget, Please install and try again."
    exit 1
fi
