# Stat Hub

A smart Hub for holding server stat

[![Build Status](https://secure.travis-ci.org/likexian/stathub-go.png)](https://secure.travis-ci.org/likexian/stathub-go)

## Overview

Stat Hub is a service for collecting and displaying servers stat.

Stat Hub have two parts, one is the SERVER for recving, storing and displaying stat, the other is the CLIENT for collecting and sending stat to SERVER. Just two binary files is needed for all of this.

## DEMO

![demo](demo.png)

## Install

### Download on the SERVER side (master)

    https://github.com/likexian/stathub-go/releases/download/v0.5.2/stathub_linux_0.5.2.tar.gz

### Untar and move

    tar zxvf stathub_*.tar.gz
    mv stathub /var
    cd /var/stathub

### Run it

If your server os is 32 bits, please run

    ./server_x86

If your server os is 64 bits, please run

    ./server_x86_64

### Open on the brower

The default url is

    http://ip:15944

Then enter the default password: likexian

### Add a CLIENT (node)

Follow the guide

    http://ip:15944/help

## LICENSE

Copyright 2015, Li Kexian

Apache License, Version 2.0

## About

- [Li Kexian](http://www.likexian.com/)
