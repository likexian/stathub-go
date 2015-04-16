# Stat Hub

A smart Hub for holding server stat

[![Build Status](https://secure.travis-ci.org/likexian/stathub-go.png)](https://secure.travis-ci.org/likexian/stathub-go)

[English README](README-ZH.md) | [中文说明](README-ZH.md)

## 总揽

Stat Hub 是一个帮您收集并展示众多服务器状态的服务。

它由两部分组成，一是服务端，用于接收、储存和展示状态；另一个是客户端，它用于收集并发送状态到服务端。而这一切，您只需要两个二进制文件。

## 演示

![demo](demo.png)

## 安装

### 在服务端下载

    https://github.com/likexian/stathub-go/releases/download/v0.5.2/stathub_linux_0.5.2.tar.gz

### 解压并移动

    tar zxvf stathub_*.tar.gz
    mv stathub /var
    cd /var/stathub

### 运行它

如果您的系统是32位的，运行

    ./server_x86

如果您的系统是64位的，运行

    ./server_x86_64

### 在您电脑的浏览器打开它

默认URL是

    http://ip:15944

输入默认密码: likexian

### 添加一个客户端

按以下提示操作

    http://ip:15944/help

## 版权

Copyright 2015, Li Kexian

Apache License, Version 2.0

## 关于

- [Li Kexian](http://www.likexian.com/)
