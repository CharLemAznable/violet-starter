#!/usr/bin/env bash

violet_version="v0.3.0"

curl -LJO https://github.com/CharLemAznable/violet-starter/releases/download/$violet_version/violet-etcd.$violet_version.linux.amd64.tar.xz

tar -xvJf violet-etcd.$violet_version.linux.amd64.tar.xz

mv violet-etcd.$violet_version.linux.amd64.bin violet-etcd
