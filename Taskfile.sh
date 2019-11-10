#!/bin/bash

function help {
    echo "$0 <task> <args>"
    echo "Tasks:"
    compgen -A function | cat -n
}

function deploy-when-master {
    if $(git branch | grep \* | cut -d ' ' -f2 | grep -q master); then
      echo "Deploying..."
      curl -X POST -F token=$INFRA_TOKEN -F "ref=master" -F "variables[SERVICE_NAME]=tasktab" -F "variables[SERVICE_VERSION]=pipeline-$CI_PIPELINE_ID" https://gitlab.com/api/v4/projects/6986946/trigger/pipeline
    fi
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