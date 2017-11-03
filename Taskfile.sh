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
    go build
}

function prod {
    TRAFFIC_ENV=production ./gotag
}

function default {
    dev
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}