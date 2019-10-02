#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function dump_schema {
  sudo docker-compose exec db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev --result-file=/dump/schema.sql"
}
function build {
    docker build -t tasktab .
}
function up {
    docker-compose up -d
}
function stop {
    docker-compose stop
}
function default {
    up
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}