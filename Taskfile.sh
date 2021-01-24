#!/bin/bash

function help() {
  echo "$0 <task> <args>"
  echo "Tasks:"
  compgen -A function | cat -n
}

function build() {
  docker build -t tasktab .
}

function dev-backend() {
  DEV_MODE=true go run . www
}

function dev-dump-schema() {
  docker-compose exec -T db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev" | grep -v "Using a password on the command line interface can be insecure" | sed 's/ AUTO_INCREMENT=[0-9]*//g' >migrations/0.sql
}

function default() {
  help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}
