#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function dump-schema {
  sudo docker-compose exec db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev" | grep -v "Using a password on the command line interface can be insecure" | sed 's/ AUTO_INCREMENT=[0-9]*//g' > migrations/0.sql
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