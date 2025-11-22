# UI Service - React Portfolio

This is the frontend UI service for the home-server project, built with React.

## Development

### Install Dependencies
```bash
make install-ui
# or
cd ui-service && npm install
```

### Run Development Server
```bash
make dev-ui
# or
cd ui-service && npm start
```

The dev server will run on `http://localhost:3000`

## Production Build

### Build for Production
```bash
make build-ui
```

This script will:
1. Build the React app (`npm run build` in ui-service)
2. Copy the build output to `gateway/ui-build`
3. The gateway will serve these static files

### How It Works

1. **Build Process**: Running `make build-ui` creates optimized production files in `ui-service/build/`
2. **Copy to Gateway**: The build files are copied to `gateway/ui-build/`
3. **Gateway Serves UI**: The gateway serves:
   - Static assets (JS, CSS, images) from `/static/*`
   - React app HTML for all non-API routes (supports client-side routing)

## Gateway Integration

The gateway (`gateway/app/main.go`) is configured to:

1. **API Routes**: All backend services under `/api/v1/*`
   - `/api/v1/auth/*` - Authentication service
   - `/api/v1/stats/*` - Stats service  
   - `/api/v1/files/*` - File service
   - `/api/v1/camera/*` - Camera service

2. **Static Files**: Served from `ui-build/static/*`
   - CSS bundles
   - JavaScript bundles
   - Media files

3. **SPA Routing**: All other routes serve `index.html` for React Router

## Architecture

```
Gateway (port 8080)
├── /api/v1/*          → Backend microservices
├── /static/*          → React static assets
├── /health            → Gateway health check
└── /*                 → React SPA (index.html)
```

## Production Deployment

1. Build the UI:
   ```bash
   make build-ui
   ```

2. Build and start all services:
   ```bash
   make up!
   ```

3. Access the application:
   - Gateway + UI: `http://localhost:8080`
   - API endpoints: `http://localhost:8080/api/v1/*`

## Notes

- The UI is served by the gateway, not as a separate service
- In production, the React app makes API calls to `/api/v1/*` (same origin, no CORS issues)
- The gateway handles authentication and proxying to backend services
- All static assets are bundled and optimized for production
