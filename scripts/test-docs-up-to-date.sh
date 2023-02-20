#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
IFS=$'\n\t'

# We first generate new docs.
make docs

# We have to update the index (since timestamps might have changed).
git update-index --refresh
# We then check to see if there are any differences.
output=$(git diff-index --exit-code --patch --stat HEAD -- docs)
up_to_date=$?
if test ${up_to_date} -ne 0; then
    echo 'Docs are not up to date'
    echo 'Please run `make docs` and commit the changes'
    echo 'Uncommitted changes:'
    echo "${output}"
    exit 1
fi
