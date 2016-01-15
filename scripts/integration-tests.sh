#!/bin/bash

base_docker_image() {
  pushd $DIR/mq-resource-tester
  docker build --tag aminjam/mq-resource-tester .
  popd
}

docker_test() {
  echo "-->$FUNCNAME $@"
  pushd $DIR
  local dockerfile=$1
  local plugin=$2
  local image="aminjam/test"
  docker build --file $dockerfile --tag ${image} .
  docker run -it --rm -e "PLUGIN=$plugin" ${image}
  popd
}
main() {
  set -eo pipefail

  readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
  local specific_plugin=$1
  local dockerfiles=$(ls $DIR/Dockerfile.test.*)

  base_docker_image

  for dockerfile in $dockerfiles; do
    local name=$(basename $dockerfile)
    local plugin="${name/Dockerfile.test./}"
    if [ -z "${specific_plugin}" ]; then
      docker_test $dockerfile $plugin
    elif [ -n "${specific_plugin}" -a "${specific_plugin}" == "$plugin" ]; then
      docker_test $dockerfile $plugin
    fi
  done
}

main "$@"
