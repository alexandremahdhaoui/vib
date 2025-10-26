#!/usr/bin/env sh

set -o nounset
set -o errexit

__VERSION="${1}"

export GOPATH="${GOPATH:-$(go env GOPATH)}"
export GOBIN="${GOBIN:-${GOPATH}/bin}"
export PATH="${GOBIN}:${PATH}"

echo "-- [INFO] installing vib"
go install "github.com/alexandremahdhaoui/vib/cmd/vib@${__VERSION}"

__HOSTNAME="$(hostname)"
# Regex the hostname, if invalid set profile name to default
if echo "${__HOSTNAME}" | grep -qE '^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$' && [ ${#__HOSTNAME} -le 63 ]; then
    VIB_PROFILE="${__HOSTNAME}"
else
    echo "-- [WARN] Hostname \"${__HOSTNAME}\" is not a valid resource name."
    while true; do
        echo -n "Please enter a profile name: "
        read VIB_PROFILE
        if echo "${VIB_PROFILE}" | grep -qE '^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$' && [ ${#VIB_PROFILE} -le 63 ]; then
            break
        else
            echo "-- [ERROR] Invalid profile name. Please use only alphanumeric characters, hyphens, and a maximum of 63 characters."
        fi
    done
fi
echo "-- [INFO] creating default profile \"${VIB_PROFILE}\""
vib create profile "${VIB_PROFILE}"

__SHELL_RC="${HOME}/.$(basename "${SHELL}")rc"
echo "-- [INFO] adding vib profile \"${VIB_PROFILE}\" to \"${__SHELL_RC}\""
cat <<EOF | tee -a "${__SHELL_RC}"
. <(vib render profile "${VIB_PROFILE}")
EOF
