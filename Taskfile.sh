#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function build {
    docker build -t tasktab .
}
function up {
    docker-compose up -d
    cd frontend
    npm run dev
}
function stop {
    docker-compose stop
}

function default {
    up
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}