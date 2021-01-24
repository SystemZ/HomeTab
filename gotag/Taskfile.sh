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
       go build
       ./gotag serve &
       # exclude all files with dots
       inotifywait -e modify -r `pwd` --exclude cache --exclude gotag.sqlite3-journal --exclude '/\..+' --exclude cache
       kill_dev
   done
}

function build-req {
    # works on ubuntu 14.04/16.04
    sudo apt-get install -y build-essential g++-arm-linux-gnueabihf mingw-w64
}

function build {
    sudo echo "Building frontend..."
    cd frontend
    yarn build
    cd ../
    echo "Building backend..."
    #go build -ldflags='-w -s -extldflags "-static"' -a -o gotag
    go build -o gotag
    sudo docker build -t gotag:latest -f Dockerfile .
    echo "Cleaning up..."
    rm gotag
    rm -r frontend/dist
}

function build_old {
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
    build-windows-amd64
}

function up {
    sudo docker-compose up
}

function build-linux-amd64 {
    echo "Building linux-amd64 ..."
    go build
    zip -r9 builds/gotag-$BUILD_VERSION_STR-linux-amd64.zip gotag migrations templates LICENSE README.md
    rm gotag
}

function build-linux-arm {
    echo "Building linux-arm ..."
    CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o gotag
    zip -r9 builds/gotag-$BUILD_VERSION_STR-linux-arm.zip gotag migrations templates LICENSE README.md
    rm gotag
}

function build-windows-amd64 {
    echo "Building windows-amd64 ..."
    CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o gotag.exe
    zip -r9 builds/gotag-$BUILD_VERSION_STR-windows-amd64.zip gotag.exe migrations templates LICENSE README.md
    rm gotag.exe
}

function prod {
    TRAFFIC_ENV=production ./gotag
}

function default {
    help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}
