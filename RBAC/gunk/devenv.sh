#!/usr/bin/env bash

readonly gunk_dir=$(git rev-parse --show-toplevel)/gunk
readonly tools_dir=$(git rev-parse --show-toplevel)/tools

PATH=$gunk_dir/bin:$gunk_dir/node_modules/.bin:$PATH

export PATH
