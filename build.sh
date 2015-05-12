#!/bin/bash

[ ! -d stathub ] && mkdir stathub
rm -rf stathub/*_$(uname -m)

cd server
go build
cd ..
mv server/server stathub/server_$(uname -m)

cd client
go build
cd ..
mv client/client stathub/client_$(uname -m)
