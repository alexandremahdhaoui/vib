#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

VIB="go run ./cmd/vib"
VIB_PATH="${HOME}/.config/vib"
TMP_FILE="$(mktemp)"

function __cleanup() {
	echo "-- Cleaning up..."
	rm "${TMP_FILE}"
	${VIB} delete expressionset default-es
	${VIB} delete -n nst-0 expressionset nst-0-es
	${VIB} delete -n nst-1 profile nst-1-p
}

trap __cleanup EXIT

echo "-- Starting smoke tests"

# -- Simple tests

echo "-- Running simple tests..."
${VIB} create profile st-0
${VIB} create profile st-1
${VIB} create profile st-2
${VIB} get profile
${VIB} delete profile st-{0,1,2}
echo "-- ✅ Simple tests pass"

# -- Namespaced tests

echo "Running namespaced tests..."
${VIB} create expressionset default-es
${VIB} create -n nst-0 expressionset nst-0-es
${VIB} create -n nst-1 profile nst-1-p

cat <<EOF | ${VIB} apply -f -
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: ExpressionSet
metadata:
  name: default-es
spec:
  arbitraryKeys:
    - key
  resolverRef:
    name: plain
    namespace: vib-system
EOF

cat <<EOF | ${VIB} apply -f -
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: ExpressionSet
metadata:
  name: nst-0-es
  namespace: nst-0
spec:
  keyValues:
    - key: value
  resolverRef:
    name: function
    namespace: vib-system
EOF

cat <<EOF | ${VIB} apply -f -
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: Profile
metadata:
  name: nst-1-p
  namespace: nst-1
spec:
  refs:
    - name: default-es
      namespace: default
    - name: nst-0-es
      namespace: nst-0
EOF

${VIB} render expressionset default-es 1>"${TMP_FILE}"
${VIB} render -n nst-0 expressionset nst-0-es 1>>"${TMP_FILE}"
${VIB} render -n nst-1 profile nst-1-p 1>>"${TMP_FILE}"

diff "${TMP_FILE}" <(
	cat <<EOF
key
function key() {
value
}
key
function key() {
value
}
EOF
)
echo "-- ✅ Namespaced tests pass"

# -- Success

echo "-- 🎉 Smoke tests ran successfully"
