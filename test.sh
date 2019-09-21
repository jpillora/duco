#!/bin/bash
set -e
echo "zipping..."
zip function.zip function.sh bootstrap
echo "deploy..."
# aws lambda create-function \
# --handler function.handler \
# --runtime provided \
# --role arn:aws:iam::652507618334:role/lambda-role
aws lambda update-function-code \
  --function-name bash-runtime \
  --zip-file fileb://function.zip
echo "invoke..."
aws lambda invoke \
  --function-name bash-runtime \
  --payload '{"text":"Hello"}' \
  /tmp/out.txt
echo "output:"
cat /tmp/out.txt
rm /tmp/out.txt
