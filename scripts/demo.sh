#!/usr/bin/env sh
set -eu
: "${INFRAI_API_KEY:?set INFRAI_API_KEY}"
curl -sS -X POST http://localhost:8080/contact -H 'Content-Type: application/json' -d '{"name":"Ada","email":"ada@example.com","message":"The CLI build log is useful."}'
