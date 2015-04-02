# Stat Hub

A smart Hub for holding server stat

[![Build Status](https://secure.travis-ci.org/likexian/stathub-go.png)](https://secure.travis-ci.org/likexian/stathub-go)

## Overview

Stat Hub is a service for collecting and displaying servers stat.

Stat Hub have two parts, one is the SERVER for recving, storing and displaying stat, the other is the CLIENT for collecting and sending stat to SERVER.

And just TWO binary files is need for all of this.

## DEMO

![demo](demo.png)

## Install

### Download on the SERVER sile (master)

- Linux 32Bit: https://github.com/likexian/stathub-go/releases/download/v0.5.1/stathub_linux_x86_0.5.1.tar.gz
- Linux 64Bit: https://github.com/likexian/stathub-go/releases/download/v0.5.1/stathub_linux_x86_64.0.5.1.tar.gz

### Untar and move

    tar zxvf stathub_*.tar.gz
    sudo mv stathub /var

### Run it

    cd /var/stathub
    ./server

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
