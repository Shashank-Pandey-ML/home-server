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

### 1. User Login

```javascript
const response = await fetch('/api/v1/auth/login', {
  method: 'POST',
  body: JSON.stringify({ email: 'user@example.com', password: 'password123' })
});

const data = await response.json();
// Returns: access_token, refresh_token, token_type, expires_in
localStorage.setItem('access_token', data.access_token);
localStorage.setItem('refresh_token', data.refresh_token);
```

### 2. Making Authenticated Requests

```javascript
const token = localStorage.getItem('access_token');
const response = await fetch('/api/v1/stats/summary', {
  headers: { 'Authorization': `Bearer ${token}` }
});
```

### 3. Token Validation Process

When a request hits a protected route:

1. **Gateway extracts token** from `Authorization: Bearer <token>` header
2. **Gateway validates locally** using cached RSA public key (~0.5ms)
3. **Gateway verifies** JWT signature, expiration, and claims
4. **Gateway proceeds** if valid, or returns 401 if invalid
5. **Gateway proxies request** to target service with user context

No network call to auth-service required!

### 4. Token Refresh

```javascript
const refreshToken = localStorage.getItem('refresh_token');
const response = await fetch('/api/v1/auth/refresh', {
  method: 'POST',
  body: JSON.stringify({ refresh_token: refreshToken })
});
const data = await response.json();
localStorage.setItem('access_token', data.access_token);
```

### 5. Logout

```javascript
await fetch('/api/v1/auth/logout', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${accessToken}` },
  body: JSON.stringify({ refresh_token: refreshToken })
});
localStorage.removeItem('access_token');
localStorage.removeItem('refresh_token');
```

## Implementation Details

### Gateway Auth Middleware

Located in: `gateway/middleware/auth.go`

**AuthMiddleware()** - Requires valid JWT token:
- Extracts token from `Authorization` header
- Validates token locally using cached RSA public key
- Sets user context (`user_id`, `email`, `is_admin`)
- Returns 401 if invalid

**OptionalAuthMiddleware()** - Token validation if present:
- Validates token if `Authorization` header exists
- Sets user context if valid
- Continues without error if no token

### Public Key Endpoint

Located in: `auth-service/handlers/handlers.go`

**GET /api/v1/auth/public-key**

Response:
```json
{
  "public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
  "algorithm": "RS256",
  "key_type": "RSA"
}
```

Gateway fetches and caches this public key for 1 hour.

## React Integration

### API Service Helper

```javascript
// src/services/api.js
class APIClient {
  getAuthHeaders() {
    const token = localStorage.getItem('access_token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
  }

  async request(endpoint, options = {}) {
    const response = await fetch(`/api/v1${endpoint}`, {
      ...options,
      headers: { ...this.getAuthHeaders(), ...options.headers }
    });
    
    if (response.status === 401) {
      const refreshed = await this.refreshToken();
      if (refreshed) return this.request(endpoint, options);
      window.location.href = '/login';
    }
    
    return response.json();
  }

  async refreshToken() {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) return false;
    
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken })
    });
    
    if (response.ok) {
      const data = await response.json();
      localStorage.setItem('access_token', data.access_token);
      return true;
    }
    return false;
  }
}
```

### Protected Route Component

```javascript
// src/components/ProtectedRoute.js
function ProtectedRoute({ children }) {
  const token = localStorage.getItem('access_token');
  return token ? children : <Navigate to="/login" replace />;
}
```

## Testing Authentication

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"email": "user@example.com", "password": "password123"}'

# 2. Protected route (no token) - Returns 401
curl http://localhost:8080/api/v1/stats/summary

# 3. Protected route (with token) - Returns data
curl http://localhost:8080/api/v1/stats/summary \
  -H "Authorization: Bearer <token>"

# 4. Public route - Works without token
curl http://localhost:8080/
```

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
