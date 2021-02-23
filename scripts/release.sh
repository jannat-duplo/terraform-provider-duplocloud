#!/bin/bash -eu
# This script starts or finishes a release using git flow
#
echo_version() {
    sed -ne 's/^VERSION=\([0-9\.]*\).*$/\1/p' <Makefile
}

release_start() {
    local version
    version="$(echo_version)"

    GIT_MERGE_AUTOEDIT=no git flow release start "$version"
}

release_finish() {
    local version
    version="$(echo_version)"

    GIT_MERGE_AUTOEDIT=no git flow release finish "$version"
    git push origin develop:develop
    git push origin "$version"
    git push origin master:master
}

release_set_version() {
    local version
    version="$1"

    trap 'rm -f Makefile~ examples/*/main.tf~' EXIT
    
    sed -e 's/^\(VERSION=\)[0-9\.]*/\1'"$version"'/g' -i~ Makefile
    find examples -name main.tf -exec \
        sed -e 's/\(version = "\)[0-9\.]*\(".*# RELEASE VERSION\)/\1'"$version"'\2/g' -i~ \{\} \;
}

action="$1" ; shift
case "$action" in
start)
    release_start "$@"
    ;;
finish)
    release_finish "$@"
    ;;
next)
    release_set_version "$1"
    GIT_MERGE_AUTOEDIT=no git commit -m '[release] bump version to '"$1"
    ;;
set-version)
    release_set_version "$1"
    ;;
*)
    echo "Usage: release.sh {start|finish|next}"
    exit 1
esac