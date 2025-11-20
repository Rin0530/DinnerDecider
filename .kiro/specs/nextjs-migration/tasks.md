# Implementation Plan: Next.js Migration

- [x] 1. Set up Next.js project structure and dependencies
  - Install Next.js and required dependencies (next, @types/node)
  - Update package.json scripts for Next.js commands (dev, build, start)
  - Create next.config.js with API proxy configuration and standalone output
  - Update tsconfig.json for Next.js compatibility
  - _Requirements: 1.1, 1.4, 4.1, 4.2, 4.3, 4.5_

- [x] 2. Create Next.js App Router structure
  - [x] 2.1 Create root layout component
    - Create app/layout.tsx with HTML structure and global providers
    - Move global styles from src/styles/globals.css to app/globals.css
    - Include Header and Footer components in layout
    - Add metadata configuration for SEO
    - _Requirements: 1.1, 1.5, 6.1, 6.3_

  - [x] 2.2 Create error handling components
    - Create app/error.tsx for error boundary functionality
    - Create app/not-found.tsx for 404 handling
    - Create app/loading.tsx with LoadingSpinner component
    - _Requirements: 2.5, 6.1_

- [x] 3. Migrate page components to Next.js
  - [x] 3.1 Migrate Home page
    - Create app/page.tsx from src/pages/Home.tsx
    - Add 'use client' directive for interactivity
    - Update imports to use Next.js path aliases
    - Add page metadata
    - _Requirements: 1.2, 1.3, 2.1, 6.1, 6.3_

  - [x] 3.2 Migrate Result page
    - Create app/result/page.tsx from src/pages/Result.tsx
    - Add 'use client' directive
    - Update imports and metadata
    - _Requirements: 1.2, 1.3, 2.1, 6.1, 6.3_

  - [x] 3.3 Migrate Meals page
    - Create app/meals/page.tsx from src/pages/Meals.tsx
    - Add 'use client' directive
    - Update imports and metadata
    - _Requirements: 1.2, 1.3, 2.1, 6.1, 6.3_

  - [x] 3.4 Migrate Settings page
    - Create app/settings/page.tsx from src/pages/Settings.tsx
    - Add 'use client' directive
    - Update imports and metadata
    - _Requirements: 1.2, 1.3, 2.1, 6.1, 6.3_

- [x] 4. Update components for Next.js compatibility
  - [x] 4.1 Update Header and Footer components
    - Replace react-router-dom Link with next/link
    - Update navigation logic for Next.js routing
    - Add 'use client' directive to Header
    - _Requirements: 2.1, 6.1_

  - [x] 4.2 Update all interactive components with 'use client' directive
    - Add 'use client' to IngredientForm, IngredientList, MealCard components
    - Ensure all components work with Next.js rendering
    - _Requirements: 2.1, 6.1_

- [x] 5. Update services and API client for Next.js
  - [x] 5.1 Update API client configuration
    - Change environment variable from import.meta.env.VITE_API_BASE_URL to process.env.NEXT_PUBLIC_API_BASE_URL
    - Ensure API client works in both client and server contexts
    - _Requirements: 3.1, 3.2_

  - [x] 5.2 Create environment variable files
    - Create .env.local with NEXT_PUBLIC_API_BASE_URL for development
    - Update .env.example with Next.js environment variables
    - _Requirements: 3.2_

- [x] 6. Update state management for Next.js
  - [x] 6.1 Update Zustand stores
    - Ensure stores work with Next.js client components
    - Add client-side checks for localStorage access in recipeStore
    - Verify store initialization works correctly
    - _Requirements: 2.2_

- [x] 7. Update testing infrastructure for Next.js




  - [x] 7.1 Create Vitest configuration for Next.js


    - Create vitest.config.ts to work with Next.js structure
    - Add path alias resolution for @/ imports
    - Ensure jsdom environment is properly configured
    - Remove test config from vite.config.ts
    - _Requirements: 5.1, 5.3_

  - [x] 7.2 Update test files


    - Update imports in all test files to use Next.js paths
    - Mock next/navigation router in tests (replace react-router-dom mocks)
    - Update environment variable mocks for Next.js
    - Update test utils to remove BrowserRouter wrapper
    - _Requirements: 5.2, 5.3_

  - [ ]* 7.3 Verify all tests pass
    - Run test suite and fix any failing tests
    - Ensure MSW integration still works
    - _Requirements: 5.2, 5.4_
-

- [x] 8. Update Docker configuration





  - [x] 8.1 Update Dockerfile for Next.js


    - Modify Dockerfile to use Next.js standalone build
    - Copy .next/standalone, .next/static, and public directories
    - Update CMD to run Next.js server (node server.js)
    - Fix Node version mismatch (use node:20-alpine for both stages)
    - _Requirements: 3.4, 4.4_

  - [ ]* 8.2 Test Docker build and deployment
    - Build Docker image and verify it works
    - Test API communication between frontend and backend containers
    - Verify environment variables are passed correctly
    - _Requirements: 3.1, 3.4, 4.4_
-

- [-] 9. Implement Next.js optimizations





  - [x]* 9.1 Add Image optimization


    - Replace img tags with next/image where appropriate
    - Configure image domains in next.config.js if needed
    - _Requirements: 6.2_

  - [x]* 9.2 Add font optimization


    - Use next/font for web font optimization
    - Update layout to include optimized fonts
    - _Requirements: 6.4_

  - [x]* 9.3 Verify code splitting


    - Ensure Next.js automatic code splitting is working
    - Verify lazy loading of components if used
    - _Requirements: 6.5_

- [x] 10. Clean up Vite-specific files and dependencies







  - [x] 10.1 Remove Vite dependencies




    - Remove vite, @vitejs/plugin-react, vite-plugin-pwa from package.json
    - Remove react-router-dom dependency (if exists)
    - Run pnpm install to update lockfile
    - _Requirements: 1.4_

  - [x] 10.2 Remove Vite configuration files


    - Delete vite.config.ts
    - Delete index.html
    - Delete src/main.tsx
    - Delete src/vite-env.d.ts
    - Delete src/App.tsx (router configuration)
    - _Requirements: 1.4_

  - [x] 10.3 Remove old directory structure


    - Remove src/pages directory after verifying all pages are migrated
    - Remove src/layouts directory
    - Remove dist directory (Vite build output)
    - _Requirements: 1.4, 6.1_

  - [ ]* 10.4 Update documentation
    - Update README.md with Next.js commands and setup instructions
    - Remove references to Vite in documentation
    - _Requirements: 4.1, 4.2, 4.3_

- [x] 11. Final verification








  - [x] 11.1 Verify development server works


    - Run `pnpm dev` and test all pages load correctly
    - Test navigation between pages
    - Verify hot module replacement works
    - _Requirements: 1.2, 1.3, 4.1_

  - [x] 11.2 Verify production build works


    - Run `pnpm build` and ensure build succeeds
    - Run `pnpm start` and test production server
    - Verify all pages work in production mode
    - _Requirements: 4.2, 4.3_

  - [ ]* 11.3 Test API integration end-to-end
    - Verify API calls work from all pages
    - Test error handling for API failures
    - Test ingredient CRUD operations
    - Test recipe suggestion flow
    - _Requirements: 3.1, 3.2, 3.3_
