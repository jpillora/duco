#!/bin/bash
set -e
echo "invoke..."
aws lambda invoke \
  --function-name go-raw-runtime \
  --payload "eyJ0ZXh0IjoiSGVsbG8ifQo=" \
  /tmp/out.txt
echo "output:"
cat /tmp/out.txt | jq
rm /tmp/out.txt
