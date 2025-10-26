#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Create a temporary directory for our mock executables
MOCK_BIN_DIR="$(mktemp -d)"
export PATH="${MOCK_BIN_DIR}:${PATH}"

# Create a dummy rc file to avoid tee error
DUMMY_RC_FILE="${HOME}/.bashrc"
touch "${DUMMY_RC_FILE}"

# Mock the hostname command to return an invalid hostname.
cat > "${MOCK_BIN_DIR}/hostname" <<EOF
#!/bin/sh
echo "invalid-hostname.with-dots"
EOF
chmod +x "${MOCK_BIN_DIR}/hostname"

# Mock the vib command to check the profile name.
cat > "${MOCK_BIN_DIR}/vib" <<EOF
#!/bin/sh
if [ "\$1" = "create" ] && [ "\$2" = "profile" ]; then
    if [ "\$3" = "my-profile" ]; then
        echo "Test passed: Profile name is 'my-profile' when user provides it."
    else
        echo "Test failed: Profile name should be 'my-profile', but got '\$3'."
        # Clean up before exiting
        rm -rf "${MOCK_BIN_DIR}"
        exit 1
    fi
# Allow other vib commands to pass without error
elif [ "\$1" = "render" ]; then
    echo "vib render command mocked"
fi
EOF
chmod +x "${MOCK_BIN_DIR}/vib"

# Mock the go command to prevent it from running.
cat > "${MOCK_BIN_DIR}/go" <<EOF
#!/bin/sh
echo "Mocked go command"
EOF
chmod +x "${MOCK_BIN_DIR}/go"

# Run the installer script, piping "my-profile" to the prompt.
echo "my-profile" | ./cmd/vib-installer/vib-installer.sh v0.0.0

# Clean up the dummy file and mock directory
rm -f "${DUMMY_RC_FILE}"
rm -rf "${MOCK_BIN_DIR}"
