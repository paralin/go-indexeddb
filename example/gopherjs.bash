#!/bin/bash
set -eo pipefail

echo "Browse to http://localhost:8080/github.com/paralin/go-indexeddb/example/"
gopherjs serve -v
