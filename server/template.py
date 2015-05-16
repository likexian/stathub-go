#!/usr/bin/env python
# -*- coding:utf-8 -*-

'''
/*
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 *
 * Copyright 2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */
'''

template = '''/*
 * A smart Hub for holding server stat
 * https://www.likexian.com/
 *
 * Copyright 2015, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package main


const (
    Template_Bootstrap = `%s`
    Template_Layout = `%s`
    Template_Index = `%s`
    Template_Login = `%s`
    Template_Help = `%s`
    Template_Node = `%s`
    Default_TLS_CERT = `%s`
    Default_TLS_KEY = `%s`
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
Template_Node = read_file('./template/node.html')
Default_TLS_CERT = read_file('./cert/cert.pem')
Default_TLS_KEY = read_file('./cert/cert.key')
template = template % (Template_Bootstrap, Template_Layout, Template_Index, Template_Login, Template_Help, Template_Node,
    Default_TLS_CERT, Default_TLS_KEY)
write_file('template.go', template)
