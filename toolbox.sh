#!/bin/bash
#
# SPDX-License-Identifier: MIT
# Copyright (c) 2025 Alejandro Visiedo <alejandro.visiedo@gmail.com>
#
# Prepare golang profile

pkgs+=(golang make bat)
source "${TOOLBOX_PROFILE_PATH}/toolbox.common.sh"

# Sometimes python tools are used in the repository
if [ -e "requirements.txt" ] || [ -e "requirements-dev.txt" ]; then
  source "${TOOLBOX_PROFILE_PATH}/toolbox.python.sh"
fi

if [ -e "tools/go.mod" ]; then
  # shellcheck disable=SC2317
  toolbox enter "${TOOLBOX}" <<EOF
make install-tools
EOF
fi

