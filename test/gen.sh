#!/bin/sh
xo --int32-type int32 --escape-all "mysql://root:@localhost/testxo?charset=utf8&parseTime=True&loc=Local" -o ./out --template-path ../templates