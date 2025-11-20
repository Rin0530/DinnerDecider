# Requirements Document

## Introduction

DinnerDecider は、ユーザが冷蔵庫の中身を登録すると、その食材を元に GPT が夕食の献立を提案する Web アプリケーションです。本ドキュメントでは、フロントエンドの要件を定義します。フロントエンドは React + TypeScript で構築され、モバイルファーストのレスポンシブデザイン、直感的な UI/UX、テスト可能で保守しやすいモジュール構成を提供します。

## Glossary

- **System**: DinnerDecider フロントエンドアプリケーション
- **User**: アプリケーションを使用するエンドユーザ
- **Ingredient**: 冷蔵庫に登録される食材
- **Recipe Suggestion**: GPT が提案する献立
- **Backend API**: フロントエンドがデータを取得・送信するバックエンドサービス
- **Store**: アプリケーションの状態を管理する Zustand ストア

## Requirements

### Requirement 1

**User Story:** As a User, I want to view all ingredients in my refrigerator, so that I can see what food items I currently have.

#### Acceptance Criteria

1. WHEN the User navigates to the Home page, THE System SHALL retrieve the list of ingredients from the Backend API using GET /ingredients
2. WHEN the Backend API returns the ingredient list, THE System SHALL display each ingredient with its name, quantity, and expiration date
3. IF the Backend API returns an error, THEN THE System SHALL display an error message to the User
4. THE System SHALL render the ingredient list within 0.5 seconds of receiving the API response

### Requirement 2

**User Story:** As a User, I want to add new ingredients to my refrigerator, so that I can keep track of what food I have available.

#### Acceptance Criteria

1. WHEN the User clicks the "Add Ingredient" button on the Home page, THE System SHALL display an ingredient input form
2. WHEN the User submits the form with valid ingredient data, THE System SHALL send a POST request to /ingredients with the ingredient information
3. WHEN the Backend API successfully creates the ingredient, THE System SHALL add the new ingredient to the displayed list
4. IF the User submits invalid data, THEN THE System SHALL display validation error messages
5. THE System SHALL support keyboard navigation for the ingredient input form

### Requirement 3

**User Story:** As a User, I want to edit existing ingredients, so that I can update quantities or expiration dates when they change.

#### Acceptance Criteria

1. WHEN the User clicks the "Edit" button on an ingredient, THE System SHALL display an edit form pre-filled with the current ingredient data
2. WHEN the User submits the edit form with valid data, THE System SHALL send a PUT request to /ingredients/:id with the updated information
3. WHEN the Backend API successfully updates the ingredient, THE System SHALL update the ingredient in the displayed list
4. IF the Backend API returns an error, THEN THE System SHALL display an error message and revert to the previous state

### Requirement 4

**User Story:** As a User, I want to delete ingredients from my refrigerator, so that I can remove items I no longer have.

#### Acceptance Criteria

1. WHEN the User clicks the "Delete" button on an ingredient, THE System SHALL display a confirmation dialog
2. WHEN the User confirms the deletion, THE System SHALL send a DELETE request to /ingredients/:id
3. WHEN the Backend API successfully deletes the ingredient, THE System SHALL remove the ingredient from the displayed list
4. IF the User cancels the deletion, THEN THE System SHALL close the confirmation dialog without making any changes

### Requirement 5

**User Story:** As a User, I want to request recipe suggestions based on my ingredients, so that I can decide what to cook for dinner.

#### Acceptance Criteria

1. WHEN the User clicks the "Suggest Recipes" button on the Home page, THE System SHALL send a POST request to /recipes/suggestion
2. WHEN the Backend API returns recipe suggestions, THE System SHALL store the suggestions in the Store
3. WHEN the suggestions are stored, THE System SHALL navigate the User to the Result page
4. THE System SHALL display a loading indicator while waiting for the Backend API response
5. IF the Backend API returns an error, THEN THE System SHALL display an error message to the User

### Requirement 6

**User Story:** As a User, I want to view the suggested recipes, so that I can choose what to cook for dinner.

#### Acceptance Criteria

1. WHEN the User navigates to the Result page, THE System SHALL retrieve recipe suggestions from the Store
2. THE System SHALL display each recipe suggestion with its name, cooking steps, and missing ingredients
3. WHERE no recipe suggestions are available, THE System SHALL display a message indicating no suggestions are available
4. THE System SHALL render each recipe using the MealCard component

### Requirement 7

**User Story:** As a User, I want to save my favorite recipes, so that I can easily access them later.

#### Acceptance Criteria

1. WHEN the User clicks the "Save to Favorites" button on a recipe, THE System SHALL store the recipe in localStorage
2. WHEN the recipe is saved, THE System SHALL display a confirmation message to the User
3. THE System SHALL persist favorite recipes across browser sessions using localStorage

### Requirement 8

**User Story:** As a User, I want the application to be accessible via keyboard, so that I can use it without a mouse.

#### Acceptance Criteria

1. THE System SHALL support keyboard navigation for all interactive elements
2. THE System SHALL provide ARIA landmarks for main content areas
3. THE System SHALL ensure all interactive elements are focusable and have visible focus indicators
4. THE System SHALL meet WCAG AA contrast requirements for all text and interactive elements

### Requirement 9

**User Story:** As a User, I want the application to work on mobile devices, so that I can use it on my phone or tablet.

#### Acceptance Criteria

1. THE System SHALL implement a responsive design that adapts to screen sizes from 320px to 1920px width
2. THE System SHALL prioritize mobile-first design patterns
3. THE System SHALL ensure touch targets are at least 44x44 pixels for mobile devices
4. THE System SHALL render correctly on iOS Safari, Android Chrome, and desktop browsers

### Requirement 10

**User Story:** As a User, I want the application to load quickly, so that I can start using it without waiting.

#### Acceptance Criteria

1. THE System SHALL complete initial rendering within 0.5 seconds on a standard broadband connection
2. THE System SHALL implement code splitting to reduce initial bundle size
3. THE System SHALL lazy load non-critical components and routes
4. THE System SHALL cache static assets using Service Worker for offline support
