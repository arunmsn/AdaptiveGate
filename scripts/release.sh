#!/usr/bin/env bash
# release.sh — cut a new ixr release
# usage: ./scripts/release.sh v0.2.0
set -euo pipefail

VERSION="${1:?usage: release.sh <version>}"

echo "Releasing $VERSION..."

# Validate semver format
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: version must be semver (e.g. v0.2.0)"
  exit 1
fi

# Ensure working tree is clean
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Error: working tree is not clean"
  exit 1
fi

git tag -s "$VERSION" -m "release $VERSION"
git push origin "$VERSION"

echo "Tag $VERSION pushed. GitHub Actions will build and publish the release."
