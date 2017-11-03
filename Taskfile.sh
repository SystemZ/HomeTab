#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function dev {
    go run *.go serve
}

function scan {
    echo $1
    go run *.go scan $1
}

function build {
    mkdir -p builds
    go build
    tar czvf builds/gotag_0.1.0_Linux-64bit.tar.gz gotag migrations templates
}

function prod {
    TRAFFIC_ENV=production ./gotag
}

function default {
    dev
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}