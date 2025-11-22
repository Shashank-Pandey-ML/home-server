#!/bin/bash

# Script to build React UI and copy to gateway

set -e

echo "Building React UI..."
cd ui-service
npm run build

echo "Copying build files to gateway..."
cd ..
rm -rf gateway/ui-build
cp -r ui-service/build gateway/ui-build

echo "✅ UI build complete and copied to gateway/ui-build"
echo "The gateway will now serve the React app from the build directory"
