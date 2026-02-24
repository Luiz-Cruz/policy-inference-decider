#!/usr/bin/env bash
set -e

LAMBDA_URL="${1:-$LAMBDA_URL}"
if [[ -z "$LAMBDA_URL" ]]; then
  echo "Usage: $0 <LAMBDA_URL> [test-file.yaml]"
  echo "Example: $0 https://xxx.lambda-url.us-east-1.on.aws"
  echo "Example: $0 https://xxx.lambda-url.us-east-1.on.aws test-100rps-mixed.yaml"
  exit 1
fi
LAMBDA_URL="${LAMBDA_URL%/}"

DIR="$(cd "$(dirname "$0")" && pwd)"
TEST_FILE="${2:-test-50rps.yaml}"
if [[ "$TEST_FILE" == */* ]]; then
  YAML="$TEST_FILE"
else
  YAML="$DIR/$TEST_FILE"
fi
if [[ ! -f "$YAML" ]]; then
  echo "Error: test file not found: $YAML"
  exit 1
fi
TMP="${TMPDIR:-/tmp}/artillery-$$.yaml"
sed "s|https://LAMBDA_FUNCTION_URL|$LAMBDA_URL|g" "$YAML" > "$TMP"
trap 'rm -f "$TMP"' EXIT

ARTILLERY_REGION="${ARTILLERY_REGION:-us-east-1}"
MEMORY="${ARTILLERY_LAMBDA_MEMORY:-3008}"
echo "Running Artillery from Lambda (region $ARTILLERY_REGION, memory ${MEMORY}MB)..."
if command -v artillery >/dev/null 2>&1; then
  artillery run-lambda --region "$ARTILLERY_REGION" --memory-size "$MEMORY" "$TMP"
else
  npx artillery run-lambda --region "$ARTILLERY_REGION" --memory-size "$MEMORY" "$TMP"
fi
