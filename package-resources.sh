#!/usr/bin/env bash

set -euo pipefail
PATH="/opt/homebrew/bin:$PATH"

# fpm-based packager for KrankyBear shared Resources
# Creates a separate package containing only the shared Resources directory
# (Images, Sounds, LICENSE) that all KrankyBear apps depend on.
#
# This package installs to /opt/local/bin/Resources/ and is a dependency
# for KrankyBearTimer, KrankyBearClock, KrankyBearNotify, etc.
#
# Benefits:
# - Shared files are owned by ONE package (no conflicts)
# - Individual apps can be installed/removed without affecting resources
# - Resources only get removed when this package is explicitly removed

usage() {
  cat <<EOF
Usage: ./package-resources.sh [linux] [ENV_VARS]

Arguments:
  linux       Build Linux packages (.deb and .rpm)

Environment variables (optional):
  VERSION     Package version (default: 1.0.0)
  ITERATION   Package iteration/release (default: 1)
  ARCH        Target arch: amd64|arm64 (default: amd64)
  OUTDIR      Output directory (default: ./installers)
  MAINTAINER  Maintainer (default: amarillier@gmail.com)
  VENDOR      Vendor (default: KrankyBear)
  URL         Project URL (default: https://github.com/amarillier/KrankyBearTimer)
  LICENSE     License (default: GNU GPL v3)

Examples:
  # Build Linux packages
  ./package-resources.sh linux
  ./package-resources.sh linux VERSION=1.0.1
  
  # Install both packages together:
  # apt install ./krankybear-resources_1.0.0-1_amd64.deb ./KrankyBearTimer_0.9.5-1_amd64.deb
  # or with dpkg:
  # dpkg -i ./krankybear-resources_1.0.0-1_amd64.deb
  # dpkg -i ./KrankyBearTimer_0.9.5-1_amd64.deb
EOF
}

# Parse command-line arguments
if [[ $# -eq 0 ]]
then
  usage
  exit 0
fi

TYPE_ARG="${1:-}"
shift  # Remove first argument, keep rest for env vars

# Validate TYPE argument
case "$TYPE_ARG" in
  linux)
    # Valid argument
    ;;
  -h|--help|-?)
    usage
    exit 0
    ;;
  *)
    echo "Error: Invalid argument '$TYPE_ARG'. Must be 'linux'." >&2
    echo "  (macOS apps include Resources in the .app bundle, no separate package needed)"
    echo ""
    usage
    exit 1
    ;;
esac

# Process remaining arguments as environment variable assignments
for arg in "$@"
do
  if [[ "$arg" == *"="* ]]
  then
    key="${arg%%=*}"
    value="${arg#*=}"
    export "$key=$value"
  fi
done

# Configurable via env vars
NAME=${NAME:-krankybear-resources}
VERSION=${VERSION:-1.0.0}
ITERATION=${ITERATION:-1}
OUTDIR=${OUTDIR:-./installers}
MAINTAINER=${MAINTAINER:-"amarillier@gmail.com"}
VENDOR=${VENDOR:-"KrankyBear"}
URL=${URL:-"https://github.com/amarillier/KrankyBearTimer"}
LICENSE=${LICENSE:-"GNU GPL v3"}

# Efficiently copy directory trees while avoiding macOS extended attributes.
copy_tree() {
  local src="$1"
  local dest="$2"
  if [[ ! -d "$src" ]]
  then
    echo "Error: Missing directory $src" >&2
    exit 1
  fi
  mkdir -p "$dest"
  if command -v rsync >/dev/null 2>&1
  then
    rsync -a --delete "$src"/ "$dest"/
  else
    (cd "$src" && tar -cf - .) | (cd "$dest" && tar -xf -)
  fi
}

# Architecture handling
ARCH="${ARCH:-amd64}"
case "$ARCH" in
  amd64|x86_64)
    DEB_ARCH=amd64
    RPM_ARCH=x86_64
    ;;
  arm64|aarch64)
    DEB_ARCH=arm64
    RPM_ARCH=aarch64
    ;;
  *)
    DEB_ARCH="$ARCH"
    RPM_ARCH="$ARCH"
    ;;
esac

# Source assets
SRC_RESOURCES="Resources"
if [[ -f "LICENSE" ]]
then
  SRC_LICENSE_FILE="LICENSE"
elif [[ -f "license.txt" ]]
then
  SRC_LICENSE_FILE="license.txt"
else
  SRC_LICENSE_FILE=""
fi

# Validate sources
if [[ ! -d "$SRC_RESOURCES" ]]
then
  echo "Error: Missing directory $SRC_RESOURCES" >&2
  exit 1
fi
if [[ ! -d "$SRC_RESOURCES/Images" ]]
then
  echo "Error: Missing directory $SRC_RESOURCES/Images" >&2
  exit 1
fi
if [[ ! -d "$SRC_RESOURCES/Sounds" ]]
then
  echo "Error: Missing directory $SRC_RESOURCES/Sounds" >&2
  exit 1
fi

mkdir -p "$OUTDIR"

# Create staging directory for Linux packages
echo "Creating staging directory for Resources package..."
STAGING_DIR=$(mktemp -d -t fpm-resources-staging.XXXXXX)
cleanup_staging() { rm -rf "$STAGING_DIR"; }
trap cleanup_staging EXIT INT TERM

# Copy the entire Resources tree
copy_tree "$SRC_RESOURCES" "$STAGING_DIR/Resources"

# Ensure LICENSE is present
if [[ -f "$STAGING_DIR/Resources/license.txt" ]] && [[ ! -f "$STAGING_DIR/Resources/LICENSE" ]]
then
  cp "$STAGING_DIR/Resources/license.txt" "$STAGING_DIR/Resources/LICENSE"
elif [[ -n "$SRC_LICENSE_FILE" ]] && [[ -f "$SRC_LICENSE_FILE" ]]
then
  cp "$SRC_LICENSE_FILE" "$STAGING_DIR/Resources/LICENSE"
fi

# Set permissions
chmod -R 755 "$STAGING_DIR/Resources"
find "$STAGING_DIR/Resources" -type f -exec chmod 644 {} \;

# Remove macOS extended attributes when building from Darwin hosts
if command -v xattr >/dev/null 2>&1
then
  xattr -cr "$STAGING_DIR" 2>/dev/null || true
  find "$STAGING_DIR" -name '._*' -delete 2>/dev/null || true
fi
if [[ "$(uname -s)" == "Darwin" ]]
then
  export COPYFILE_DISABLE=1
fi

# Common fpm arguments
COMMON_ARGS=(
  -s dir
  -n "$NAME"
  -v "$VERSION"
  --iteration "$ITERATION"
  --maintainer "$MAINTAINER"
  --vendor "$VENDOR"
  --url "$URL"
  --license "$LICENSE"
  --description "KrankyBear shared resources - Images, Sounds, and LICENSE files shared by all KrankyBear applications"
  -f
)

# Check for fpm
if ! command -v fpm >/dev/null 2>&1
then
  echo "Error: fpm not found. Install with: gem install fpm" >&2
  exit 1
fi

# Build .deb package
DEB_OUTFILE="$OUTDIR/krankybear-resources_${VERSION}-${ITERATION}_${DEB_ARCH}.deb"
echo "Building .deb ($DEB_ARCH) -> $DEB_OUTFILE..."

DEB_ARGS=(
  "${COMMON_ARGS[@]}"
  -t deb
  -a "$DEB_ARCH"
  --deb-no-default-config-files
  --directories /opt/local/bin
  --directories /opt/local/bin/Resources
)

# Mark LICENSE as config file so it's preserved
DEB_ARGS+=(--config-files "/opt/local/bin/Resources/LICENSE")

fpm \
  "${DEB_ARGS[@]}" \
  --package "$DEB_OUTFILE" \
  "$STAGING_DIR/Resources=/opt/local/bin"

echo ""

# Build .rpm package
RPM_OUTFILE="$OUTDIR/krankybear-resources_${VERSION}-${ITERATION}_${RPM_ARCH}.rpm"
echo "Building .rpm ($RPM_ARCH) -> $RPM_OUTFILE..."

RPM_ARGS=(
  "${COMMON_ARGS[@]}"
  -t rpm
  -a "$RPM_ARCH"
  --rpm-os linux
  --rpm-auto-add-directories
)

# Mark LICENSE as config file
RPM_ARGS+=(--config-files "/opt/local/bin/Resources/LICENSE")

fpm \
  "${RPM_ARGS[@]}" \
  --package "$RPM_OUTFILE" \
  "$STAGING_DIR/Resources=/opt/local/bin"

echo ""
echo "Done. Packages created:"
echo "  $DEB_OUTFILE"
echo "  $RPM_OUTFILE"
echo ""
echo "To install apps with their resources dependency:"
echo "  apt install ./$DEB_OUTFILE ./KrankyBearTimer_VERSION_amd64.deb"
echo "  # or with dpkg (install resources first):"
echo "  dpkg -i $DEB_OUTFILE"
echo "  dpkg -i ./KrankyBearTimer_VERSION_amd64.deb"

# "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942

