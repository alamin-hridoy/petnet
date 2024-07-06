#!/usr/bin/env bash
set -e

# fail on OS X old bash version (for **)
shopt -s globstar
set -euo pipefail

./install_deps.sh

source ./devenv.sh

find ./drp/v1 -mindepth 1 -maxdepth 1 -type d -execdir bash -c 'cd "$1"; gunk generate' _ {} \;
find ./dsa/v1 -mindepth 1 -maxdepth 1 -type d -execdir bash -c 'cd "$1"; gunk generate' _ {} \;
find ./dsa/v2 -mindepth 1 -maxdepth 1 -type d -execdir bash -c 'cd "$1"; gunk generate' _ {} \;
find ./v1 -mindepth 1 -maxdepth 1 -type d -execdir bash -c 'cd "$1"; gunk generate' _ {} \;

swagger-combine config_v1.json -o swaggerv1.json
swagger2openapi -p swaggerv1.json -o apiv1.yaml
rm swaggerv1.json

# swagger-combine config_v2.json -o swaggerv2tmp.json
# jq -s '.[0] * .[1]' dev_servers.json swaggerv2tmp.json >swaggerdev.json
# jq -s '.[0] * .[1]' prod_servers.json swaggerv2tmp.json >swaggerprod.json
# swagger2openapi -p swaggerdev.json -o apiv2.yaml
# swagger2openapi -p swaggerprod.json -o apiv2prod.yaml
# rm swaggerv2tmp.json swaggerdev.json swaggerprod.json
