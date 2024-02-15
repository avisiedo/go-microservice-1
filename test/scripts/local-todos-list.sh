#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

unset CREDS

exec "${REPOBASEDIR}/scripts/curl.sh" -i "${BASE_URL}/todos"
