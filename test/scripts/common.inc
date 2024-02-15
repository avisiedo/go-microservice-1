#
# Include file with common parts shared for local and ephemeral
#

# Troubleshooting:
# - Run with DEBUG=1 to see some traces from curl.sh wrapper
#   $ DEBUG=1 ./test/scripts/local-todos-create.sh
# - Run with more verbose by:
#   $ DEBUG=1 bash -xv ./test/scripts/local-todos-create.sh

function error {
    local err=$?
    printf "ERROR: %s\n" "$*" >&2
    exit $err
}

SRCDIR="$(dirname "${BASH_SOURCE[0]}")"
# shellcheck disable=SC2034  # ignore unused variable
BASEDIR="$(dirname "$(dirname "${SRCDIR}")")"
REPOBASEDIR="$(git rev-parse --show-toplevel)"
export REPOBASEDIR
