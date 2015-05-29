# Stat Hub

A smart Hub for holding server stat

[![Build Status](https://secure.travis-ci.org/likexian/stathub-go.png)](https://secure.travis-ci.org/likexian/stathub-go)

[中文说明](README-ZH.md) | [English README](README.md)

## Overview

Stat Hub is a service for collecting and displaying servers stat.

Stat Hub have two parts, one is the SERVER for recving, storing and displaying stat, the other is the CLIENT for collecting and sending stat to SERVER. Just two binary files is needed for all of this.

## DEMO

![demo](demo.png)

## Feature

- Powered by Go lang
- Only two binary files, doing and displaying data
- Easy deploy, no depends, no database required
- SSL support, your domain support, secure and easy

## Install

You shall choose a server for master, and install it following the below

### Linux with curl

    curl --insecure https://github.com/likexian/stathub-go/blob/master/setup.sh | sh

### Linux with wget

    wget --no-check-certificate -O - https://github.com/likexian/stathub-go/blob/master/setup.sh | sh

### Open on your PC browser

For most systems, it will install and start after this, now you can view it via your PC browser.

The default url is

    https://ip:15944

Then enter the default password: likexian

### Add a CLIENT (node)

Follow the guide

    https://ip:15944/help

## FAQ

- Why no data on the page ?

    Please add a client first, the guide will show how to.

- A client added, still no list on the page ?

    Please refer to the client.log on the node, it may show what is wrong.

- Shall I run the client on the server side ?

    Sure, server is a client on the other side.

- Can I use my domain instead of the https://ip:15944 ?

    Yes, please add an A record on your dns server for one subdomain, then using https://subdomain.yourdomain:15944 .

- Can I use https for visist ?

    Sure, it is SSL by default, using self-signed cert.

- Cal I use my SSL cert, not the self-signed one ?

    Yes, please replace the cert in the cert dir, and it strong recommend that doing this.

- Can I deploy it with my nginx ?

    Sure, please add the folling config to nginx conf file.

    location /stathub/ {
        proxy_pass https://127.0.0.1:15944;
        proxy_set_header X-Real-IP $remote_addr;
    }

## LICENSE

Copyright 2015, Li Kexian

Apache License, Version 2.0

## About

- [Li Kexian](https://www.likexian.com/)

## Thanks

- [All stargazers](https://github.com/likexian/stathub-go/stargazers)
- [Livid](https://github.com/livid)
- [Septembers](https://github.com/Septembers)
- [Bluek404](https://github.com/Bluek404)
- [renjie45](https://github.com/renjie45)
- [vijaygadde](https://github.com/vijaygadde)
- [davoola](https://github.com/davoola)
