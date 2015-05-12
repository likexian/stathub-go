#!/bin/bash

pname=$(uname -m)
if [ "$pname" != "x86_64" ]; then
    pname="x86"
fi

[ ! -d stathub ] && mkdir stathub
rm -rf stathub/*_$pname

cd server
go build
cd ..
mv server/server stathub/server_$pname

cd client
go build
cd ..
mv client/client stathub/client_$pname
