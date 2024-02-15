#!/bin/bash

##
# curl helper to wrap and add headers automatically based
# into the environment variables defined; the idea is
# define once, use many times, so the curl command can be
# simplified.
#
# NOTE: Do not forget to to 'export MYVAR' for the next
#       one to get a better user experience, else you need
#       to set before the command in the same command,
#       reducing the user experience then.
#
# REQUEST_ID if it is empty, a random value
#
##

# Uncomment to print verbose traces into stderr
# DEBUG=1
function verbose {
    [ "${DEBUG}" != 1 ] && return 0
    echo "$@" >&2
}

# Initialize the array of options
opts=()

# Optionally add CREDS if it was set (used for testing in ephemeral)
# See: make ephemeral-namespace-describe
if [ "${CREDS}" != "" ]; then
    opts+=("-u" "${CREDS}")
    # shellcheck disable=SC2016
    verbose '-u "${CREDS}"'
fi

# Generate a REQUEST_ID if it is not set
if [ "${REQUEST_ID}" == "" ]; then
    if [ "$(uname -s)" == "Linux" ]; then
        REQUEST_ID="test_$(cat "/proc/sys/kernel/random/uuid" | head -c 20 | base64 -w0)"
    elif [ "$(uname -s)" == "Darwin" ]; then
        REQUEST_ID="test_$(cat "/dev/random" | head -c 20 | base64)"
    else
        error "No support for REQUEST_ID on '$(uname -s)' system"
    fi
fi
opts+=("-H" "Request-Id: ${REQUEST_ID}")
verbose "-H Request-Id: ${REQUEST_ID}"

# Add Content-Type
opts+=("-H" "Content-Type: application/json")
verbose "-H Content-Type: application/json"

# Add the rest of values
opts+=("$@")

verbose /usr/bin/curl "${opts[@]}"
/usr/bin/curl "${opts[@]}"
