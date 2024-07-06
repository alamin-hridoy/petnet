#!/usr/bin/env bash

set -euo pipefail

source ./devenv.sh

pushd "$tools_dir"

GOBIN=$gunk_dir/bin go install \
	github.com/gunk/gunk 

popd

pushd "$gunk_dir"
pnpm install --frozen-lockfile
popd
