#!/bin/bash

set +xae

service=$1
cmd=$2
[[ -z $service ]] && echo 'arg 1 (service) not set' && exit 1
[[ -z $cmd ]] && echo 'arg 2 (cmd) not set' && exit 1

air --build.cmd "go build -o ./.tmp/$service-$cmd ./services/$service/cmd/$cmd" --build.bin ".tmp/$service-$cmd"
