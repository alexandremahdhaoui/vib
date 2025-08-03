#!/usr/bin/env sh

set -o nounset
set -o errexit

__VERSION="${1}"

export GOPATH="${GOPATH:-$(go env GOPATH)}"
export GOBIN="${GOBIN:-${GOPATH}/bin}"
export PATH="${GOBIN}:${PATH}"

echo "-- [INFO] installing vib"
go install "github.com/alexandremahdhaoui/vib/cmd/vib@${__VERSION}"

VIB_PROFILE="$(hostname)"
echo "-- [INFO] creating default profile \"${VIB_PROFILE}\""
vib create profile "${VIB_PROFILE}"

echo "-- [INFO] adding vib profile \"${VIB_PROFILE}\" to \".${SHELL}rc\""
cat <<EOF | tee -a "${HOME}/.${SHELL}rc"
. <(vib render profile "${VIB_PROFILE}")
EOF
