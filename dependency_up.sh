#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Get the directory of this script
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$script_dir/docker"

compose_files=( -f "docker-compose.yml" )

if [[ $# -gt 0 && "$*" == *"--dev"* ]]; then
    compose_files+=( -f "docker-compose.dev.yml")
fi

docker compose "${compose_files[@]}" up -d
