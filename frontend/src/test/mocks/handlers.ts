import { http, HttpResponse } from 'msw';
import type { Ingredient, RecipeSuggestionResponse } from '../../types';

const API_BASE_URL = 'http://localhost:8080';

// Mock data
const mockIngredients: Ingredient[] = [
  { id: 1, name: 'にんじん', quantity: '2本', expiration_date: '2025-11-01' },
  { id: 2, name: 'じゃがいも', quantity: '3個', expiration_date: '2025-11-05' },
];

const mockRecipes: RecipeSuggestionResponse = {
  suggestions: [
    {
      name: 'カレーライス',
      steps: ['野菜を切る', 'カレールーを入れる', '煮込む'],
      missing_items: ['カレールー', '肉'],
    },
    {
      name: '肉じゃが',
      steps: ['材料を切る', '煮る'],
      missing_items: ['肉', '醤油'],
    },
  ],
};

export const handlers = [
  // GET /ingredients
  http.get(`${API_BASE_URL}/ingredients`, () => {
    return HttpResponse.json(mockIngredients);
  }),

  // POST /ingredients
  http.post(`${API_BASE_URL}/ingredients`, async ({ request }) => {
    const body = (await request.json()) as Ingredient;
    const newIngredient: Ingredient = {
      ...body,
      id: Date.now(),
    };
    return HttpResponse.json(newIngredient, { status: 201 });
  }),

  // PUT /ingredients/:id
  http.put(`${API_BASE_URL}/ingredients/:id`, async ({ params, request }) => {
    const { id } = params;
    const body = (await request.json()) as Partial<Ingredient>;
    const updatedIngredient: Ingredient = {
      id: Number(id),
      name: body.name || 'Updated',
      quantity: body.quantity,
      expiration_date: body.expiration_date,
    };
    return HttpResponse.json(updatedIngredient);
  }),

  // DELETE /ingredients/:id
  http.delete(`${API_BASE_URL}/ingredients/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // POST /recipes/suggestion
  http.post(`${API_BASE_URL}/recipes/suggestion`, () => {
    return HttpResponse.json(mockRecipes);
  }),
];
