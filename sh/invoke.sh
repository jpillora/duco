#!/bin/bash
set -e
echo "invoke..."
aws lambda invoke \
  --function-name bash-runtime \
  --payload '{"text":"Hello"}' \
  /tmp/out.txt
echo "output:"
cat /tmp/out.txt
rm /tmp/out.txt
