#!/bin/bash
BUILD_VERSION_STR=""

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
       inotifywait -e modify -r `pwd` --exclude cache --exclude gotag.sqlite3-journal --exclude '/\..+' --exclude cache
       kill_dev
   done
}

function scan {
    echo $1
    go run *.go scan $1
}

function build-req {
    # works on ubuntu 14.04/16.04
    apt-get install -y build-essential g++-arm-linux-gnueabihf
}

function build {
    [ -z "$1" ] && echo "Provide commit or version string" >&2 && exit 1
    mkdir -p builds

    # if we have two parameters, use second as a version
    if [ -z "$2" ]
    then
      BUILD_VERSION_STR=$1
    else
      BUILD_VERSION_STR=$2
    fi

    build-linux-amd64
    build-linux-arm
}

function build-linux-amd64 {
    echo "Building linux-amd64 ..."
    go build
    zip -r9 builds/gotag-$BUILD_VERSION_STR-linux-amd64.zip gotag migrations templates LICENSE README.md
}

function build-linux-arm {
    echo "Building linux-arm ..."
    CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o gotag
    zip -r9 builds/gotag-$BUILD_VERSION_STR-linux-arm.zip gotag migrations templates LICENSE README.md
}

function prod {
    TRAFFIC_ENV=production ./gotag
}

function default {
    serve
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}