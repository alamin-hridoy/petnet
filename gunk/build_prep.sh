#!/bin/bash

set -ex

if [[ $drone_branch == "production" ]]; then
	cp gunk/apiv1prod.yaml docs/api.yaml
else
	cp gunk/apiv1.yaml docs/api.yaml
fi
