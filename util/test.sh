#!/usr/bin/env bash

go test -v -count 2
# go test -test.bench=".*"

# 当前目录下生成 cpu.profile 和 util.test 文件
go test -run=NONE -bench=. -cpuprofile=cpu.out -v
go tool pprof -text -nodecount=5 ./util.test cpu.out

#go test -cpuprofile=cpu.out
#go test -memprofile=mem.out
#go test -coverprofile=c.out
#go tool cover -html=c.out     #查看报告
#go tool pprof mem.out