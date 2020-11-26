#!/bin/bash
set -e
echo "invoke..."
aws lambda invoke \
  --function-name bash-runtime \
  --payload "eyJ0ZXh0IjoiSGVsbG8ifQo=" \
  --log-type "Tail" \
  /tmp/out.txt | jq -r '.LogResult' | base64 -d
echo "output:"
cat /tmp/out.txt
rm /tmp/out.txt
