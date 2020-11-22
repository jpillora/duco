#!/bin/bash
set -e
echo "invoke..."
aws lambda invoke \
  --function-name bash-runtime \
  --payload "eyJ0ZXh0IjoiSGVsbG8ifQo=" \
  /tmp/out.txt
echo "output:"
cat /tmp/out.txt
rm /tmp/out.txt
