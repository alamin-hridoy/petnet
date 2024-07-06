#!/usr/bin/env bash

set -euo pipefail

source ./devenv.sh

pushd "$tools_dir"

GOBIN=$gunk_dir/bin go install -mod=readonly github.com/gunk/gunk
GOBIN=$gunk_dir/bin go install github.com/gunk/scopegen@v0.1.1

pushd "$tools_dir/pgunk"
GOBIN=$gunk_dir/bin go install
popd

popd

pushd "$gunk_dir"
pnpm install --frozen-lockfile
popd
