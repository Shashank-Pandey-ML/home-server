#!/bin/bash

# Script to create a user in the auth database
# Usage: ./scripts/create_user.sh

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Create User ===${NC}"
echo ""

# Prompt for user details
read -p "Email: " EMAIL
read -sp "Password: " PASSWORD
echo ""
read -p "Name: " NAME
read -p "Is Admin? (y/n): " IS_ADMIN_INPUT

# Convert to boolean
if [[ "$IS_ADMIN_INPUT" =~ ^[Yy]$ ]]; then
    IS_ADMIN="true"
else
    IS_ADMIN="false"
fi

echo ""
echo -e "${BLUE}Creating user...${NC}"

# Hash password using bcrypt (cost 10)
# Build and run the Go hash utility
cd "$(dirname "$0")"
PASSWORD_HASH=$(go run hash_password.go "$PASSWORD")
cd - > /dev/null

if [ -z "$PASSWORD_HASH" ]; then
    echo -e "${RED}❌ Failed to hash password${NC}"
    exit 1
fi

# Insert user into database
docker exec -i postgres psql -U postgres -d auth << EOF
INSERT INTO users (email, password, name, is_admin, created_at, updated_at)
VALUES ('$EMAIL', '$PASSWORD_HASH', '$NAME', $IS_ADMIN, NOW(), NOW())
ON CONFLICT (email) DO UPDATE 
SET password = EXCLUDED.password,
    name = EXCLUDED.name,
    is_admin = EXCLUDED.is_admin,
    updated_at = NOW();
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ User created successfully!${NC}"
    echo ""
    echo "Email: $EMAIL"
    echo "Name: $NAME"
    echo "Admin: $IS_ADMIN"
else
    echo -e "${RED}❌ Failed to create user${NC}"
    exit 1
fi
