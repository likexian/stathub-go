#!/bin/sh

pname=$(uname -m)
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

cp CHANGS.md stathub
cp CHANGS-ZH.md stathub
cp LICENSE.md stathub
cp README.md stathub
cp README-ZH.md stathub
