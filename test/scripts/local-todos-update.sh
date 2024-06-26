#!/bin/bash
set -eo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/local.inc"

UUID="$1"
[ "${UUID}" != "" ] || error "UUID is empty"

exec "${REPOBASEDIR}/test/scripts/curl.sh" -i -X PUT -d @<(sed -e 's/{{subscription_manager_id}}/6f324116-b3d2-11ed-8a37-482ae3863d30/g' < "${REPOBASEDIR}/test/data/http/update-rhel-idm-domain.json") "${BASE_URL}/todos/${UUID}"
