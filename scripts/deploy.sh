#!/usr/bin/env bash 
set -e -o pipefail

trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd "$(dirname "$0")" && cd ..
CURRENT=$(pwd)

function release
{    
  go mod vendor
  sudo rm -rf "$CURRENT"/dist "$CURRENT"/gopath  
  export GOPATH="$CURRENT"/gopath

  tag=$1
  if [ -z "$tag" ]; then
    echo "not found tag name"
    exit 1
  fi
 
  git tag -a "$tag" -m "Add $tag"
  git push origin "$tag"

  goreleaser release --clean
}

CMD=$1
shift
"$CMD" "$*"
