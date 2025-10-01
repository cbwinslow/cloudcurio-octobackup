import { validateEnv } from '@/lib/env-validator';

// Run validation when the module is imported
const { success, error } = validateEnv();

if (!success) {
  console.error('❌ Environment validation failed:');
  console.error(error);
  console.error('Please check your .env.local file');
  process.exit(1);
}

console.log('✅ Environment variables validated successfully');