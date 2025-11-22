# Gateway Routing Strategy

## Overview

The gateway uses a clean separation between frontend (React UI) and backend (API) routes:

```
┌─────────────────────────────────────────┐
│           Gateway (Port 8080)           │
├─────────────────────────────────────────┤
│  /                    → React UI        │
│  /profile, /about     → React Router    │
│  /api/v1/auth/*       → auth-service    │
│  /api/v1/files/*      → file-service    │
│  /api/v1/stats/*      → stats-service   │
│  /health              → Gateway health  │
└─────────────────────────────────────────┘
```

## Routing Rules

### 1. **Backend APIs** - All under `/api/v1`

```
/api/v1/auth/*      → auth-service:8080
/api/v1/files/*     → file-service:8080
/api/v1/stats/*     → stats-service:8080
/api/v1/camera/*    → camera-service:8080
```

**Benefits:**
- ✅ Consistent with auth-service pattern
- ✅ Easy to add versioning (v2, v3)
- ✅ Clear separation from UI routes
- ✅ Industry standard convention

### 2. **Frontend UI** - All other routes

```
/                   → React app (index.html)
/profile            → React Router handles
/about              → React Router handles
/files              → React Router handles
/login              → React Router handles
```

**How it works:**
- Gateway serves React's `index.html` for all non-API routes
- React Router takes over client-side routing
- Supports SPA (Single Page Application) architecture

### 3. **Special Routes**

```
/health             → Gateway health check (no /api prefix)
```

## Development Setup

### Option 1: React Dev Server (Development)

When developing, React runs on its own dev server (port 3000):

```yaml
# docker-compose.yml
services:
  ui-service:
    build: ./ui-service
    ports:
      - "3000:3000"
    volumes:
      - ./ui-service:/app
      - /app/node_modules
    command: npm start
```

**In handlers.go:** Uncomment the proxy section:
```go
func ServeReactApp() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Proxy to React dev server
        services.ProxyRequest("ui-service", c)
        return
    }
}
```

**Workflow:**
1. Run `make up` (starts all services including React dev server)
2. Gateway proxies UI requests to React dev server on port 3000
3. Hot reload works automatically
4. Access everything through `http://localhost:8080`

### Option 2: Production Build (Production/Testing)

For production, serve React's built static files:

```bash
# Build React app
cd ui-service
npm run build

# Copy build to gateway
cp -r build ../gateway/ui-build
```

**In handlers.go:** Use the static file serving (already configured):
```go
func ServeReactApp() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Serve from ./ui-build directory
        // (current implementation)
    }
}
```

**Workflow:**
1. Build React app: `npm run build`
2. Copy build to gateway's `ui-build/` folder
3. Gateway serves static files directly
4. Faster, no React dev server needed

## React Configuration

### Update React API calls to use `/api/v1` prefix:

```javascript
// src/services/api.js
const API_BASE_URL = '/api/v1';

export const authService = {
  login: async (credentials) => {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
    });
    return response.json();
  },
  
  register: async (userData) => {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(userData),
    });
    return response.json();
  },
};

export const fileService = {
  upload: async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${API_BASE_URL}/files/upload`, {
      method: 'POST',
      body: formData,
    });
    return response.json();
  },
  
  list: async () => {
    const response = await fetch(`${API_BASE_URL}/files`);
    return response.json();
  },
};
```

### React Router Setup:

```javascript
// src/App.js
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/files" element={<FileManager />} />
        <Route path="/login" element={<Login />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </Router>
  );
}
```

## Docker Compose Setup

### Complete docker-compose.yml example:

```yaml
services:
  gateway:
    build:
      context: .
      dockerfile: ./gateway/Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      - auth-service
      - ui-service
    environment:
      - ENVIRONMENT=development
    volumes:
      - ./gateway/config.yaml:/app/config.yaml
      - ./ui-service/build:/app/ui-build  # For production build

  auth-service:
    build:
      context: .
      dockerfile: ./auth-service/Dockerfile
    depends_on:
      - postgres
    env_file:
      - ./auth-service/.env

  ui-service:
    build: ./ui-service
    # Development mode
    ports:
      - "3000:3000"
    volumes:
      - ./ui-service:/app
      - /app/node_modules
    command: npm start
    
  postgres:
    image: postgres:15
    env_file:
      - ./postgres/.env
    volumes:
      - postgres-data:/var/lib/postgresql/data
```

## Request Flow Examples

### Example 1: User visits homepage

```
1. Browser → GET http://localhost:8080/
2. Gateway → Checks: Not /api/v1/* → Serve React app
3. Gateway → Returns index.html
4. React loads → React Router renders Home component
```

### Example 2: User logs in

```
1. React → POST http://localhost:8080/api/v1/auth/login
2. Gateway → Matches /api/v1/auth/* → Proxy to auth-service
3. auth-service → Validates credentials, returns JWT
4. Gateway → Returns response to React
5. React → Stores JWT, redirects to /profile
```

### Example 3: User navigates to profile

```
1. React Router → Changes URL to /profile (client-side)
2. No server request (SPA navigation)
3. React → Renders Profile component
4. Profile → Fetches data: GET /api/v1/auth/profile
5. Gateway → Proxies to auth-service
6. Response → Displayed in UI
```

## Benefits of This Approach

✅ **Clean separation**: Frontend and backend clearly separated
✅ **Standard convention**: `/api/*` is industry standard
✅ **Versioning support**: Easy to add v2, v3 APIs
✅ **SPA-friendly**: React Router works seamlessly
✅ **Development friendly**: Hot reload in dev, optimized in production
✅ **Scalable**: Easy to add new microservices
✅ **CORS-free**: Everything served from same origin

## Troubleshooting

### Issue: React routing doesn't work (404 on refresh)

**Solution:** Make sure gateway serves `index.html` for non-API routes. The current implementation handles this.

### Issue: API calls fail with CORS errors

**Solution:** CORS is already configured in middleware. If issues persist, check that requests go through gateway (not directly to services).

### Issue: Static files (CSS, JS) not loading

**Solution:** 
1. Check that `ui-build/` directory exists and contains React build
2. Verify file permissions
3. Check gateway logs for file serving errors

## Next Steps

1. ✅ Gateway routing is configured
2. ⏳ Build React app with routes
3. ⏳ Implement file upload UI
4. ⏳ Add authentication UI (login/register)
5. ⏳ Create portfolio pages
6. ⏳ Deploy to home server
