#!/usr/bin/env bash

set -o errexit
set -o pipefail

VIB="go run ./cmd/vib"

# -- Basic tests
echo Run basic tests
${VIB} create profile test0
${VIB} create profile test1
${VIB} create profile test2
${VIB} get profile
${VIB} delete profile test{0,1,2}

# -- Namespaced tests
${VIB} create expressionset default-es
${VIB} create -n test-0 expressionset test-0-es
${VIB} create -n test-1 profile test-1-p

cat <<EOF >TODO_DEFAULT_ES_FILE
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: ExpressionSet
metadata:
  name: default-es
spec:
  arbitraryKeys: null
  keyValues: null
  resolverRef:
    name: function
    namespace: vib-system
EOF

cat <<EOF >TODO_TEST_0_ES_FILE
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: ExpressionSet
metadata:
  name: test-0-es
spec:
  arbitraryKeys: null
  keyValues: null
  resolverRef:
    name: function
    namespace: vib-system
EOF

cat <<EOF >TODO_P_FILE
apiVersion: vib.alexandre.mahdhaoui.com/v1alpha1
kind: Profile
metadata:
  name: test-1-p
  namespace: test-1
spec:
  refs:
    - name: test-0-es
	  namespace: test-0
	- name: default-es
EOF

${VIB} render -n test-0 expressionset es
${VIB} render -n test-1 profile p

${VIB} delete -n test-0 expressionset es
${VIB} delete -n test-1 profile p

# -- Success
echo Smoke tests ran successfully
