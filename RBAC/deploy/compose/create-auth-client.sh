#!/bin/sh

PAYLOAD='{
   "client_id": "rbac-dev",
   "client_secret": "secret",
   "grant_types": ["authorization_code", "refresh_token", "client_credentials"],
   "response_types": ["code", "id_token"],
   "post_logout_redirect_uris": [
     "http://192.168.2.100:8080",
     "http://127.0.0.1:8080",
     "http://127.0.0.1:3003",
     "http://172.16.10.133:8080"
   ],
   "redirect_uris": [
     "http://192.168.2.100:8080/oauth2/callback",
     "http://127.0.0.1:8888/oauth2/callback",
     "http://127.0.0.1:8080/oidc-callback",
     "http://127.0.0.1:3003/oauth2/callback",
     "http://172.16.10.133:8080/oauth2/callback"
   ],
   "scope": "https://rbac.brank.as/read https://rbac.brank.as/write openid offline_access"}'

until $(curl --output /dev/null --fail --silent --show-error --location \
	--request POST http://hydra:4445/clients \
	--header 'Content-Type: application/json' \
	--header 'Accept: application/json' \
	--data-raw "$PAYLOAD" ||
	curl --output /dev/null --fail --silent --show-error --location \
		--request PUT http://hydra:4445/clients/rbac-dev \
		--header 'Content-Type: application/json' \
		--header 'Accept: application/json' \
		--data-raw "$PAYLOAD"); do
	echo "."
	sleep 2
done
