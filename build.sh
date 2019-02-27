#!/bin/bash

command_exists() {
    type "$1" &> /dev/null
}

echo '+ building stathub start'
rm -rf stathub.*.tar.gz

git add .
sed -ie "s/DEBUG\([^\=]*\)= true/DEBUG\1= false/g" src/const.go
rm -rf src/const.goe

for i in "i686" "x86_64"; do
    rm -rf tmp && mkdir tmp

    cp LICENSE tmp
    cp README.md tmp
    cp README-ZH.md tmp

    echo "building the $i stathub"
    cd src
    make $i
    if [ $? -ne 0 ]; then
        echo "building the $i stathub failed"
        echo '+ building stathub failed'
        exit 1
    fi
    cd ..
    mv bin/stathub tmp/stathub
    cp service.sh tmp/service

    echo "doing upx to binary files"
    cd tmp
    if command_exists upx; then
        upx stathub >/dev/null 2>&1
    fi

    echo "packaging the $i stathub"
    tar zcf stathub.$i.tar.gz *

    cd ..
    cp tmp/stathub.$i.tar.gz .
    rm -rf tmp
done

git checkout .

VERSION=`grep -C1 'func Version' src/version.go | grep 'return ' | awk -F'"' '{print $2}'`
sed -ie "s/VERSION=\".*\"/VERSION=\"$VERSION\"/g" setup.sh
rm -rf setup.she

echo '+ building stathub done'
