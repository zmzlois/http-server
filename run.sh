#!/bin/sh

set -euo pipefail

tmpFile=${mktemp}
go build -o "$tmpFile" app/*.go

exec "$tmpFile" "$@"
