#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Function to check Zitadel availability
check_zitadel_availability() {
    if curl -s -o /dev/null http://localhost:8080/ui/console > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Get the directory of this script
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$script_dir/docker"

compose_files=( -f "docker-compose.yml" )

if [[ $# -gt 0 && "$*" == *"--dev"* ]]; then
    compose_files+=( -f "docker-compose.dev.yml")
fi

docker compose "${compose_files[@]}" up -d

# Cannot seem to be able to use healthcheck in zitadel service in docker compose file because
# the container has no shell or cmds available like curl or wget - so unable to use `docker compose up -d --wait`
attempts=0
max_attempts=5

# Check availability before waiting
if check_zitadel_availability; then
    echo "Zitadel is already available: http://localhost:8080"
    exit 0
fi

echo "Waiting for Zitadel to become available. Sleeping for 20 seconds before checking."
sleep 20
while [[ $attempts -lt $max_attempts ]]; do
    if check_zitadel_availability; then
        echo "Zitadel is now available: http://localhost:8080"
        exit 0
    else
        ((attempts++))
        echo "Waiting for Zitadel to become available. Attempt $attempts of $max_attempts. Sleeping 5 seconds."
        sleep 5
    fi
done

echo "Timeout: Zitadel did not become available within the specified time."
exit 1
