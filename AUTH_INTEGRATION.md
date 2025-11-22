# Gateway ↔ Auth-Service Integration

## Overview

The gateway enforces authentication for protected routes by validating JWT tokens locally using the public key from auth-service.

## Architecture

```
┌─────────────┐          ┌─────────────┐          ┌─────────────┐
│   Browser   │          │   Gateway   │          │Auth-Service │
└─────────────┘          └─────────────┘          └─────────────┘
      │                         │                         │
      │  1. POST /api/v1/auth/login                     │
      │───────────────────────>│                         │
      │                        │ 2. Proxy request        │
      │                        │────────────────────────>│
      │                        │                         │
      │                        │ 3. JWT tokens returned  │
      │                        │<────────────────────────│
      │  4. Login response     │                         │
      │<───────────────────────│                         │
      │                        │                         │
      │  5. GET /api/v1/stats (with Bearer token)      │
      │───────────────────────>│                         │
      │                        │ 6. Validate token       │
      │                        │    locally (no network) │
      │                        │                         │
      │                        │ 7. Proxy to stats       │
      │                        │────────────> Stats      │
      │  8. Stats data         │              Service    │
      │<───────────────────────│                         │
```

## Routes Configuration

### Public Routes (No Auth Required)

| Route | Description | Backend |
|-------|-------------|---------|
| `GET /` | Portfolio home page | React UI |
| `GET /profile` | Profile page (public) | React UI |
| `GET /health` | Gateway health check | Gateway |
| `POST /api/v1/auth/login` | User login | auth-service |
| `POST /api/v1/auth/register` | User registration | auth-service |
| `POST /api/v1/auth/refresh` | Refresh access token | auth-service |
| `GET /api/v1/auth/public-key` | Get JWT public key | auth-service |

### Protected Routes (Auth Required)

| Route | Description | Backend |
|-------|-------------|---------|
| `POST /api/v1/auth/logout` | User logout | auth-service |
| `GET /api/v1/users/profile` | Get user profile | auth-service |
| `PUT /api/v1/users/profile` | Update user profile | auth-service |
| `ANY /api/v1/stats/*` | Stats service | stats-service |
| `ANY /api/v1/files/*` | File operations | file-service |
| `ANY /api/v1/camera/*` | Camera feeds | camera-service |
| All UI routes except `/profile` | React pages | React UI |

## Authentication Flow

### Login Process
1. User submits credentials to `/api/v1/auth/login`
2. Auth service validates and returns JWT tokens (access + refresh)
3. Client stores tokens locally

### Protected Request Flow
1. Client sends request with `Authorization: Bearer <token>` header
2. Gateway validates token locally using cached RSA public key (~0.5ms)
3. Gateway verifies JWT signature, expiration, and claims
4. If valid, gateway proxies request to target service with user context
5. If invalid, returns 401 Unauthorized

**Key benefit:** No network call to auth-service for validation!

### Token Refresh
- When access token expires, client uses refresh token
- Endpoint: `POST /api/v1/auth/refresh`
- Returns new access token

### Logout
- Invalidates refresh token in database
- Endpoint: `POST /api/v1/auth/logout`
- Client clears stored tokens

## Implementation Details

### Gateway Components
- **Location:** `gateway/middleware/auth.go`
- **AuthMiddleware:** Validates JWT token, sets user context, returns 401 if invalid
- **OptionalAuthMiddleware:** Validates token if present, continues without error if missing

### Auth Service Components
- **Location:** `auth/handlers/handlers.go`
- **Public Key Endpoint:** `GET /api/v1/auth/public-key`
  - Returns RSA public key in PEM format
  - Gateway caches for 1 hour
  - Used for local JWT validation

### React Integration Pattern
- Store tokens in localStorage (or httpOnly cookies for production)
- Add `Authorization: Bearer <token>` header to all API requests
- Implement automatic token refresh on 401 responses
- Protected routes check for token presence before rendering

## Security Features

### ✅ Implemented
- JWT tokens with RSA-256 signatures
- Local token validation (fast & scalable)
- Token expiration (configurable)
- Refresh tokens for re-authentication
- CORS, rate limiting, security headers

### 🔒 Best Practices
- Use HTTPS in production
- Short access token lifetime (15-60 min)
- Rotate refresh tokens
- Use httpOnly cookies instead of localStorage
- Implement token blacklist for immediate revocation

## Next Steps

1. ✅ Local JWT validation with public key
2. ✅ Protected routes configured
3. ⏳ Build React login/register UI
4. ⏳ Add user registration endpoint
5. ⏳ Add httpOnly cookie support
6. ⏳ Role-based access control (RBAC)
7. ⏳ Token blacklist for immediate revocation
