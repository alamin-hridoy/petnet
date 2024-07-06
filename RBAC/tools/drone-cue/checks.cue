package drone

import (
	"strings"
	"list"
)

#GoVersion:        "1.14.7"
#ProtoVersion:     "3.12.4"
#NodeMajorVersion: "14"
#NodeVersion:      #NodeMajorVersion + ".8.0"
#NodeProtoVersion: #NodeMajorVersion + ".8.0"

// Enforce our workspace to be consistent with the cache directory.
workspace: path: "/root"
#GoCacheDir:   "/go-cache"
#PnpmCacheDir: "/pnpm-cache"
#SassBinariesCacheDir: "/node-sass-binaries-cache"

// An image should be pinned to a specific version, or be an official plugin.
#Image:           (#OfficialPlugin | #UnofficialImage)
#UnofficialImage: strings.Contains(":") & !~":latest$"
#OfficialPlugin:  =~"^plugins/[a-z]+"
#Step: {
	name:  #Identifier
	image: #Image

	// Ensure consistency between Go steps.
	if strings.Contains(image, "golang") || strings.Contains(image, "ob-go-test") {
		environment: GOPATH: #GoCacheDir + "/gopath"
		volumes: [{
			name: "go-cache"
			path: #GoCacheDir
		}]
	}

	// Ensure consistency between Node steps.
	if strings.Contains(image, "node") {
		if !strings.Contains(image, "proto") && name != "js-build-bankboard" {
			image: "asia.gcr.io/b-api-production/node-deps:" + #NodeVersion
			volumes: [{
				name: "pnpm-cache"
				path: #PnpmCacheDir
			},
			{
				name: "node-sass-binaries-cache"
				path: #SassBinariesCacheDir
			}]
			environment: {
				PNPM_CACHE_FOLDER:  #PnpmCacheDir
			}
		}

		if name == "js-build-bankboard" {
			image: "asia.gcr.io/b-api-production/node-deps:" + #NodeVersion
			volumes: [{
				name: "pnpm-cache"
				path: #PnpmCacheDir
			},
			{
				name: "node-sass-binaries-cache"
				path: #SassBinariesCacheDir
			}]
			environment: {
				PNPM_CACHE_FOLDER:  #PnpmCacheDir
			}
		}

		if name == "gunk-generate" {
			image: "asia.gcr.io/b-api-production/node-deps-proto:" + #NodeProtoVersion + "-" + #ProtoVersion + "-" + #GoVersion
			volumes: [{
				name: "pnpm-cache"
				path: #PnpmCacheDir
			}, {
				name: "go-cache"
				path: #GoCacheDir
			}]
			environment: {
				PNPM_CACHE_FOLDER:  #PnpmCacheDir
			}
		}
	}

	// Ensure we use consistent Go-based images.
	if strings.Contains(image, "golang") {
		image: "golang:" + #GoVersion
	}
	if strings.Contains(image, "ob-go-test") {
		image: "asia.gcr.io/b-api-production/ob-go-test:go" + #GoVersion + "-chrome84.0.4147.105"
	}

	// Run the module tests with CGO, and with access to postgres.
	if strings.HasPrefix(name, "test-module") {
		environment: {
			DATABASE_CONNECTION: "user=postgres host=postgres port=5432 dbname=postgres sslmode=disable"
			CGO_ENABLED:         "1"
		}
	}

	// Go build steps need to depend on set-version.
	if strings.HasPrefix(name, "go-build") {
		depends_on: list.Contains("set-version")
	}
}
#Volume: {
	name: #Identifier
	host?: {
		path: string
	}

	if name == "pnpm-cache" {
    host: {
      path: "/var/lib/drone-cache/${DRONE_REPO_NAME}/pnpm/node" + #NodeMajorVersion
    }
  }
}

// Ensure step names are unique via a struct with key-value fields like a map.
for i, step in steps {
	stepsByName: {"\(step.name)": step}
}

// Ensure all steps without dependencies are at the beginning.
for i, step in steps if len(step.depends_on) > 0 {
	// TODO(mvdan): can we avoid the quadratic cost?
	for j, step2 in steps[i+1:] {
		stepsByName: {"\(step2.name)": {
			// We could use [_, ...] here, but that gives worse
			// error messages.
			depends_on: list.MinItems(1)
		}}
	}
}
