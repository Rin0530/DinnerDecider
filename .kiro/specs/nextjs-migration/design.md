# Design Document: Next.js Migration

## Overview

This document outlines the technical design for migrating the Dinner Decider frontend from a Vite-based React SPA to a Next.js application using the App Router. The migration will enable server-side rendering, improve SEO, and eliminate the dependency on index.html while preserving all existing functionality.

### Migration Strategy

The migration will follow an incremental approach:
1. Set up Next.js project structure alongside existing code
2. Migrate routing and pages to Next.js App Router
3. Migrate components, services, and state management
4. Update configuration and build process
5. Remove Vite-specific files
6. Update Docker deployment

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Next.js Application                   │
├─────────────────────────────────────────────────────────┤
│  App Router (app/)                                       │
│  ├── layout.tsx (Root Layout)                           │
│  ├── page.tsx (Home)                                    │
│  ├── result/page.tsx                                    │
│  ├── meals/page.tsx                                     │
│  └── settings/page.tsx                                  │
├─────────────────────────────────────────────────────────┤
│  Components (components/)                                │
│  ├── Header, Footer, LoadingSpinner                     │
│  ├── IngredientForm, IngredientList                     │
│  └── MealCard, ErrorBoundary                            │
├─────────────────────────────────────────────────────────┤
│  State Management (stores/)                              │
│  ├── ingredientStore.ts (Zustand)                       │
│  └── recipeStore.ts (Zustand)                           │
├─────────────────────────────────────────────────────────┤
│  Services (services/)                                    │
│  ├── api.ts (API Client)                                │
│  ├── ingredientApi.ts                                   │
│  └── recipeApi.ts                                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
                 Backend API (Port 8080)
```

### Directory Structure

```
frontend/
├── app/                          # Next.js App Router
│   ├── layout.tsx               # Root layout with providers
│   ├── page.tsx                 # Home page (/)
│   ├── result/
│   │   └── page.tsx            # Result page
│   ├── meals/
│   │   └── page.tsx            # Meals page
│   ├── settings/
│   │   └── page.tsx            # Settings page
│   └── globals.css             # Global styles
├── components/                  # React components (unchanged)
├── stores/                      # Zustand stores (unchanged)
├── services/                    # API services (updated for Next.js)
├── hooks/                       # Custom hooks (unchanged)
├── types/                       # TypeScript types (unchanged)
├── utils/                       # Utility functions (unchanged)
├── public/                      # Static assets
├── next.config.js              # Next.js configuration
├── tsconfig.json               # TypeScript configuration
└── package.json                # Dependencies
```

## Components and Interfaces

### 1. Root Layout Component

**File**: `app/layout.tsx`

```typescript
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>
        <Providers>
          <Header />
          <main>{children}</main>
          <Footer />
        </Providers>
      </body>
    </html>
  )
}
```

**Responsibilities**:
- Define HTML structure
- Include global providers (ErrorBoundary, etc.)
- Include Header and Footer components
- Load global styles

### 2. Page Components

Each page will be converted to a Next.js page component:

**File**: `app/page.tsx` (Home)
- Migrated from `src/pages/Home.tsx`
- Client component for interactivity
- Uses existing IngredientForm and IngredientList components

**File**: `app/result/page.tsx`
- Migrated from `src/pages/Result.tsx`
- Displays recipe suggestions
- Uses existing MealCard component

**File**: `app/meals/page.tsx`
- Migrated from `src/pages/Meals.tsx`
- Shows favorite recipes
- Uses Zustand store for state

**File**: `app/settings/page.tsx`
- Migrated from `src/pages/Settings.tsx`
- Application settings

### 3. Client Components

All interactive components will use `'use client'` directive:
- IngredientForm
- IngredientList
- MealCard
- LoadingSpinner
- ErrorBoundary (converted to error.tsx)

### 4. API Service Updates

**File**: `services/api.ts`

Update environment variable access:
```typescript
// Before (Vite)
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

// After (Next.js)
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL
```

### 5. State Management

Zustand stores remain unchanged but need client-side initialization:
- ingredientStore.ts
- recipeStore.ts

LocalStorage access in recipeStore will be wrapped in client-side checks.

## Data Models

No changes to existing data models:
- Ingredient
- IngredientInput
- Recipe
- RecipeSuggestionResponse
- ApiError
- ApiResponse

## Configuration

### 1. Next.js Configuration

**File**: `next.config.js`

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable React strict mode
  reactStrictMode: true,
  
  // Output standalone for Docker
  output: 'standalone',
  
  // API proxy for backend
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.BACKEND_URL || 'http://localhost:8080'}/api/:path*`,
      },
    ]
  },
  
  // Image optimization
  images: {
    domains: [],
  },
}

module.exports = nextConfig
```

### 2. Environment Variables

**File**: `.env.local` (development)
```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
BACKEND_URL=http://backend:8080
```

**File**: `.env.production` (production)
```
NEXT_PUBLIC_API_BASE_URL=/api
BACKEND_URL=http://backend:8080
```

### 3. TypeScript Configuration

Update `tsconfig.json` for Next.js:
```json
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### 4. Package.json Updates

```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "biome check .",
    "format": "biome format --write .",
    "test": "vitest --run",
    "test:watch": "vitest",
    "type-check": "tsc --noEmit"
  },
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "zustand": "^5.0.0"
  },
  "devDependencies": {
    "@biomejs/biome": "^1.9.4",
    "@testing-library/jest-dom": "^6.9.1",
    "@testing-library/react": "^16.0.1",
    "@types/node": "^20.0.0",
    "@types/react": "^18.3.11",
    "@types/react-dom": "^18.3.1",
    "autoprefixer": "^10.4.20",
    "jsdom": "^25.0.1",
    "msw": "^2.11.6",
    "postcss": "^8.4.47",
    "tailwindcss": "^3.4.14",
    "typescript": "~5.6.2",
    "vitest": "^2.1.3"
  }
}
```

**Removed dependencies**:
- vite
- @vitejs/plugin-react
- vite-plugin-pwa (will be replaced with next-pwa if needed)

**Added dependencies**:
- next
- @types/node

## Error Handling

### 1. Error Boundaries

Convert ErrorBoundary component to Next.js error handling:

**File**: `app/error.tsx`
```typescript
'use client'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <div>
      <h2>Something went wrong!</h2>
      <button onClick={() => reset()}>Try again</button>
    </div>
  )
}
```

### 2. Not Found Handling

**File**: `app/not-found.tsx`
```typescript
export default function NotFound() {
  return (
    <div>
      <h2>404 - Page Not Found</h2>
    </div>
  )
}
```

### 3. Loading States

**File**: `app/loading.tsx`
```typescript
import LoadingSpinner from '@/components/LoadingSpinner'

export default function Loading() {
  return <LoadingSpinner size="lg" message="Loading..." />
}
```

## Testing Strategy

### 1. Component Testing

Vitest configuration for Next.js:

**File**: `vitest.config.ts`
```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './'),
    },
  },
})
```

### 2. Test Updates

- Update imports to use Next.js paths
- Mock Next.js router using `next/navigation`
- Update environment variable mocks

### 3. MSW Integration

MSW will continue to work with Next.js for API mocking in tests.

## Docker Deployment

### Updated Dockerfile

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN pnpm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy standalone build
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "server.js"]
```

## Migration Steps Summary

1. **Install Next.js dependencies**
   - Add next, @types/node
   - Remove vite, @vitejs/plugin-react

2. **Create Next.js structure**
   - Create app/ directory
   - Create layout.tsx, page.tsx files

3. **Migrate pages**
   - Convert each page to Next.js page component
   - Add 'use client' directive where needed

4. **Update services**
   - Change environment variable access
   - Update API client for Next.js

5. **Update configuration**
   - Create next.config.js
   - Update tsconfig.json
   - Update package.json scripts

6. **Migrate tests**
   - Update vitest configuration
   - Update test imports and mocks

7. **Update Docker**
   - Modify Dockerfile for Next.js standalone build

8. **Remove Vite files**
   - Delete vite.config.ts
   - Delete index.html
   - Delete src/main.tsx

## SEO and Metadata

Each page will include metadata:

```typescript
export const metadata = {
  title: 'Dinner Decider - Home',
  description: 'Decide what to cook based on your ingredients',
}
```

## Performance Considerations

- Use Next.js automatic code splitting
- Implement dynamic imports for heavy components
- Use Next.js Image component for optimized images
- Enable React strict mode
- Use standalone output for smaller Docker images

## Backward Compatibility

- All existing URLs will continue to work
- API endpoints remain unchanged
- State management behavior is preserved
- Component functionality is unchanged
