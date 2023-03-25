#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
IFS=$'\n\t'

# Generate the next version.
readonly next_version=$(semver patch)

# Generate a tag.
git tag "${next_version}"

# Push the tag
git push origin "${next_version}"
