#!/bin/bash
set -eux -o pipefail

FULL_PATH="$1"
PATH_WITHOUT_PREFIX="${FULL_PATH#services/}"

SERVICE_NAME=$(echo "$PATH_WITHOUT_PREFIX" | cut -d'/' -f1)
CMD_NAME=$(echo "$PATH_WITHOUT_PREFIX" | cut -d'/' -f3)

# Export the variables
export service="$SERVICE_NAME"
export cmd="$CMD_NAME"
