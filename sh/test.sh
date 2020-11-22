#!/bin/bash
set -e
echo "zipping..."
zip function.zip bootstrap
echo "deploy..."
# aws lambda create-function \
# --handler function.handler \
# --runtime provided \
# --role arn:aws:iam::652507618334:role/lambda-role
aws lambda update-function-code \
  --function-name bash-runtime \
  --zip-file fileb://function.zip
#
rm function.zip
#test it
./invoke.sh
