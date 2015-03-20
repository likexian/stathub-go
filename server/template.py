#!/usr/bin/env python
# -*- coding:utf-8 -*-


template = '''package main


const (
    Template_Bootstrap = `%s`
    Template_Layout = `%s`
    Template_Index = `%s`
    Template_Login = `%s`
    Template_Help = `%s`
)
'''


def read_file(fname):
    fp = open(fname, 'r')
    text = fp.read()
    fp.close()
    return text


def write_file(fname, text):
    fp = open(fname, 'w')
    fp.write(text)
    fp.close()


Template_Bootstrap = read_file('./static/bootstrap.css')
Template_Layout = read_file('./template/layout.html')
Template_Index = read_file('./template/index.html')
Template_Login = read_file('./template/login.html')
Template_Help = read_file('./template/help.html')
template = template % (Template_Bootstrap, Template_Layout, Template_Index, Template_Login, Template_Help)
write_file('template.go', template)
