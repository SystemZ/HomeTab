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

function dev-seed() {
  # TODO import DB schema
  go run . user-create --username dev --password dev --email example@example.com
}

function dev-dump-schema() {
  docker-compose exec -T db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev" | grep -v "Using a password on the command line interface can be insecure" | sed 's/ AUTO_INCREMENT=[0-9]*//g' >migrations/0.sql
}

function deploy-docker() {
  [ -z "$1" ] && echo "Host not provided" >&2 && exit 1
  rm -rf new
  cd frontend || exit
  yarn build
  mv dist ../new
  cd ../
  CGO_ENABLED=0 go build -o tasktab
  docker build -t tasktab . && docker save tasktab | bzip2 | pv | ssh $1 'bunzip2 | docker load'
  # TODO figure out fully automatic and comfortable UnRAID deployment
  # /usr/local/emhttp/plugins/dynamix.docker.manager/scripts/docker run -d --name='tasktab' --net='bridge' -e TZ="Europe/Warsaw" -e HOST_OS="Unraid" -p '1337:80/tcp' 'tasktab'
}

function ci-build-backend() {
  go test -race $(go list ./... | grep -v /vendor/ | grep -v "/tasktab/agent")
  CGO_ENABLED=0 go build -o $CI_PROJECT_DIR/tasktab
  ls -alh $CI_PROJECT_DIR/
}

function ci-build-frontend() {
  cd frontend || exit 1
  apk add yarn
  yarn install --frozen-lockfile
  yarn build
  ls -alh dist
  mv dist $CI_PROJECT_DIR/new
}

function ci-build-img() {
  docker login -u gitlab-ci-token -p $CI_JOB_TOKEN $CI_REGISTRY
  docker build --tag $CI_REGISTRY_IMAGE:pipeline-$CI_PIPELINE_ID --tag $CI_REGISTRY_IMAGE:latest .
  docker push $CI_REGISTRY_IMAGE:pipeline-$CI_PIPELINE_ID
}

function default() {
  help
}

TIMEFORMAT="Task completed in %3lR"
time ${@:-default}
