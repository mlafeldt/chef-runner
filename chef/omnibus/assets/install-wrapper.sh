#!/bin/sh
# A smart wrapper around Omnibus Installer

set -e

script=$1
version=$2

manifest=/opt/chef/version-manifest.txt
current_version=$(head -n1 "$manifest" 2>/dev/null | cut -d" " -f2)

case "$version" in
""|false)
    echo "==> Doing nothing."
    ;;
latest)
    echo "==> Installing latest version of Chef..."
    sh "$script"
    ;;
true)
    if test -n "$current_version"; then
        echo "==> Chef version $current_version installed. Doing nothing."
    else
        echo "==> Installing latest version of Chef..."
        sh "$script"
    fi
    ;;
*)
    if test "$current_version" = "$version"; then
        echo "==> Chef version $version already installed. Doing nothing."
    else
        echo "==> Installing Chef version $version ..."
        sh "$script" -v "$version"
    fi
    ;;
esac
