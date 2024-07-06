#!/usr/bin/env bash
set -e

# fail on OS X old bash version (for **)
shopt -s globstar

set -euo pipefail

./install_deps.sh

source ./devenv.sh

find ./v1 -mindepth 1 -maxdepth 1 -type d -execdir bash -c 'cd "$1"; gunk generate' _ {} \;

swagger-combine config.json -o swaggerv1.json
swagger2openapi -p swaggerv1.json -o apiv1.yaml
rm swaggerv1.json

for f in v*/**/*.{js,d.ts}; do
	mkdir -p dist/$(dirname $f)
	mv $f dist/$f

	# fix the import paths
	perl -pi -e "s#([\.\.\/])+(.)*/all_pb#./all_pb#" dist/$f

	# strip unnecessary google API annotation imports
	perl -pi -e 's/.*google_api_annotations_pb.*\n//' dist/$f
	perl -pi -e 's/import.*google_api_annotations_pb.*\n//' dist/$f
done
