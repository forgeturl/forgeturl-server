#!/bin/bash

set -ve
cd $(dirname $0)

cd gensql
go run gorm_gen.go