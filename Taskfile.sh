#!/bin/bash
REPO="github.com/systemz/hometab"
IMG_NAME="hometab"

function help() {
  echo "$0 <task> <args>"
  echo "Tasks:"
  compgen -A function | cat -n
}

function install-tools() {
  go get -u github.com/go-bindata/go-bindata/...
}

function generate-resources() {
  # generate files to include in binary
  ## migrations
  go-bindata -pkg resources -o internal/resources/resources.go internal/model/migrations/
  # for recursive use dir/...
  #go-bindata data/...
}

function dev-backend() {
  generate-resources
  # build and run binary
  DEV_MODE=true TEMPLATE_PATH="$(pwd)/web/templates/" ASSETS_PATH="$(pwd)/" go run github.com/systemz/hometab/cmd/hometab serve
}

function dev-frontend() {
  # shellcheck disable=SC2164
  cd frontend
  yarn install --frozen-lockfile
  yarn serve
}

function dev-seed() {
  # TODO import DB schema
  go run . user-create --username dev --password dev --email example@example.com
}

function dev-dump-schema() {
  docker-compose exec -T db /bin/sh -c "/usr/bin/mysqldump -udev -pdev --no-data dev" | grep -v "Using a password on the command line interface can be insecure" | sed 's/ AUTO_INCREMENT=[0-9]*//g' >model/migrations/0.sql
}

# build

function build-backend() {
  generate-resources
  # shellcheck disable=SC2046
  go test -race -cover -mod=readonly $(go list ./... | grep -v hometab-agent) || exit 1
  (
    cd cmd/hometab
    CGO_ENABLED=1 go build || exit 1
  )
}

function build-frontend() {
  (
    cd frontend
    yarn install --frozen-lockfile
    yarn build
  )
}

function build-img() {
  # make disposable dir with all files necessary
  # it makes context sent to docker daemon smaller, resulting in faster builds
  # you can also retain older builds if you need it, can be useful in rapid development
  BUILD_TMP_DIR="builds/$(date +"%s")"
  mkdir -p "$BUILD_TMP_DIR"
  # backend
  cp cmd/hometab/hometab "$BUILD_TMP_DIR/hometab" || exit 1
  cp -r web/templates "$BUILD_TMP_DIR/templates" || exit 1
  # frontend
  cp -r frontend/dist "$BUILD_TMP_DIR/new" || exit 1
  # docker
  cp build/Dockerfile "$BUILD_TMP_DIR/Dockerfile" || exit 1
  # finally... build this image image already!
  docker build -t $IMG_NAME "$BUILD_TMP_DIR"
}

function build-img-clear() {
  rm -rf builds
}

# deploy

function deploy-docker() {
  [ -z "$1" ] && echo "Host not provided" >&2 && exit 1
  rm -rf new
  cd frontend || exit
  yarn build
  mv dist ../new
  cd ../
  CGO_ENABLED=0 go build -o tasktab
  docker build -t $IMG_NAME . && docker save tasktab | bzip2 | pv | ssh $1 'bunzip2 | docker load'
  # TODO figure out fully automatic and comfortable UnRAID deployment
  # /usr/local/emhttp/plugins/dynamix.docker.manager/scripts/docker run -d --name='tasktab' --net='bridge' -e TZ="Europe/Warsaw" -e HOST_OS="Unraid" -p '1337:80/tcp' 'tasktab'
}

# TODO rewrite CI stuff for github

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
