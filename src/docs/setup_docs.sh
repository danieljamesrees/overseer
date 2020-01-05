#!/usr/bin/env bash

# Only works on MacOS.

set -o errexit
set -o nounset
set -o pipefail

brew install hugo
hugo check --gc
hugo server --minify
