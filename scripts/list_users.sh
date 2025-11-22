#!/bin/bash

# Script to list all users in the database
# Usage: ./scripts/list_users.sh

set -e

# Colors for output
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Users in Database ===${NC}"
echo ""

docker exec -i postgres psql -U postgres -d auth << 'EOF'
\x auto
SELECT id, email, name, is_admin, created_at, updated_at 
FROM users 
ORDER BY created_at DESC;
EOF
