#!/bin/bash

# 查找所有.go文件
for file in $(find . -name "*.go"); do
  # 获取文件所在目录的路径
  dir=$(dirname $file)
  # 切换到该目录
  echo $dir
  cd $dir
  # 执行 go build 命令
  go build
  # 返回上一级目录
  cd -
done