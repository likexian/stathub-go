#!/bin/bash

command_exists() {
    type "$1" &> /dev/null
}

echo '+ building stathub start'

for i in "i686" "x86_64"; do
    rm -rf tmp && mkdir tmp

    cp LICENSE tmp
    cp README.md tmp
    cp README-ZH.md tmp
    cp src/VERSION tmp

    echo "building the $i stathub"
    cd src
    make $i
    cd ..
    mv bin/stathub tmp/stathub
    cp service.sh tmp/service

    echo "goupxing binary files"
    cd tmp
    if command_exists goupx; then
        goupx stathub >/dev/null 2>&1
    fi

    echo "packaging the $i stathub"
    tar zcf stathub.$i.tar.gz *

    cd ..
    cp tmp/stathub.$i.tar.gz .
    rm -rf tmp
done

echo '+ building stathub done'
