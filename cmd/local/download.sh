#!/usr/bin/env bash

violet_version="v0.1.0"

curl -LJO https://github.com/CharLemAznable/violet-starter/releases/download/$violet_version/violet-local.$violet_version.linux.amd64.tar.xz

tar -xvJf violet-local.$violet_version.linux.amd64.tar.xz

mv violet-local.$violet_version.linux.amd64.bin violet-local
