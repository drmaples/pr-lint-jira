import { defineConfig, globalIgnores } from 'eslint/config';
import tseslint from 'typescript-eslint';

export default defineConfig(
  globalIgnores(['dist']),
  {
  files: ['**/*.{ts,tsx}'],
  extends: [
      tseslint.configs.recommended,
  ],
  languageOptions: {
    ecmaVersion: 2020, // Supports modern JS features
    sourceType: 'module', // Supports ES module imports
  },
  rules: {
    'no-console': 'warn', // Warn on console logs (avoid in production code)
    'no-debugger': 'warn', // Warn on debugger statements
    'no-unused-vars': 'off', // Disable built-in no-unused-vars
    '@typescript-eslint/no-unused-vars': ['error'], // Enforce no unused vars in TypeScript
    '@typescript-eslint/no-explicit-any': 'error', // Disallow 'any' type
    '@typescript-eslint/explicit-module-boundary-types': 'error', // Enforce return types on functions
    '@typescript-eslint/no-empty-function': 'error', // Disallow empty functions
    '@typescript-eslint/explicit-function-return-type': ['error', { allowExpressions: true }], // Enforce function return types
    '@typescript-eslint/ban-ts-comment': 'error', // Disallow using // @ts-ignore and other ts-comments
    '@typescript-eslint/consistent-type-definitions': ['error', 'interface'], // Prefer interfaces over type aliases
    '@typescript-eslint/member-ordering': 'error', // Enforce consistent member ordering
    '@typescript-eslint/no-inferrable-types': 'error', // Prevent explicit types when they can be inferred
    '@typescript-eslint/no-var-requires': 'error', // Disallow require() calls
    'semi': ['error', 'never'], // Enforce semicolons at the end of statements
    'quotes': ['error', 'single'], // Enforce single quotes
  },
});
