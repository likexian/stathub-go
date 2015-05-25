#!/bin/sh

echo '+ building stathub start'
[ ! -d stathub ] && mkdir stathub
rm -rf tmp && mkdir tmp

echo 'building the server'
cd server
go build
cd ..
mv server/server tmp
cp server/service tmp

echo 'building the client'
cd client
go build
cd ..
mv client/client tmp

cp CHANGS.md tmp
cp CHANGS-ZH.md tmp
cp LICENSE.md tmp
cp README.md tmp
cp README-ZH.md tmp
cp VERSION tmp

echo 'goupxing binary files'
cd tmp
if [ -f $(which goupx) ]; then
    goupx server client >/dev/null 2>&1
fi

echo 'packaging the server'
tar zcf server_$(uname -m).tar.gz server service VERSION *.md
echo 'packaging the client'
tar zcf client_$(uname -m).tar.gz client VERSION *.md

cd ..
cp tmp/server_$(uname -m).tar.gz tmp/client_$(uname -m).tar.gz stathub
rm -rf tmp

echo '+ building stathub done'
