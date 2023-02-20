#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail
IFS=$'\n\t'

# We first generate new docs.
make docs

# We then check to see if there are any differences.
if test -n "$(git status --porcelain -- docs)"; then
    echo
    echo 'Uncommitted changes:'
    git diff -- docs
    echo
    echo 'Docs are not up to date'
    echo 'Please run `make docs` and commit the changes'
    echo
    exit 1
fi
