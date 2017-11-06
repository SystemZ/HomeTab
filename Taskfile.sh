#!/bin/bash

function kill_dev() {
    PID1=$(ps aux | grep '/tmp/go' | awk '{print $2}')
    PID2=$(ps aux | grep 'go run' | awk '{print $2}')
    kill $PID1
    kill $PID2
}

function sigint_handler() {
    kill_dev
    exit
}

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function serve {
   trap sigint_handler SIGINT
   while true; do
       go run *.go serve &
       # exclude all files with dots
       inotifywait -e modify -e move -e create -e delete -e attrib -r `pwd` --exclude cache --exclude index.lock --exclude '/\..+' --exclude cache
       kill_dev
   done
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
    serve
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}