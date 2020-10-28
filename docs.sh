#!/usr/bin/env bash

set -e

log() {
	echo "---" "$@"
}

# put the godocdown template in a temporary directory
TEMPLATE=$(mktemp)
trap 'rm ${TEMPLATE}' EXIT

cat <<EOF >"${TEMPLATE}"
# package {{ .Name }}

\`import "{{ .ImportPath }}"\`

<p>
  <a href="https://pkg.go.dev/{{ .ImportPath }}"><img src="https://img.shields.io/badge/doc-reference-007d9b?logo=go&style=flat-square" alt="go.dev" /></a>
  <a href="https://goreportcard.com/report/{{ .ImportPath }}"><img src="https://goreportcard.com/badge/{{ .ImportPath }}?style=flat-square" alt="Go Report Card" /></a>
  <a href="https://sourcegraph.com/{{ .ImportPath }}?badge"><img src="https://sourcegraph.com/{{ .ImportPath }}/-/badge.svg?style=flat-square" alt="SourceGraph" /></a>
</p>

{{ .EmitSynopsis }}

{{ .EmitUsage }}
EOF

# build the godocdown tool
GODOCDOWN=$(
	SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
	cd "${SCRIPTDIR}"
	cd "$(pwd -P)"

	IMPORT=github.com/robertkrimen/godocdown/godocdown
	go install -v "${IMPORT}"
	go list -f '{{ .Target }}' "${IMPORT}"
)

"${GODOCDOWN}" -template "${TEMPLATE}" . > README.md
