import postgres from 'https://deno.land/x/postgresjs@v3.4.5/mod.js';

export const APP_NAME = Deno.env.get('APP_NAME')!;
export const APP_VERSION = '0.1.0';
export const SERVER_PORT = parseInt(Deno.env.get('SERVER_PORT')!);
export const NUM_WORKERS = parseInt(Deno.env.get('NUM_WORKERS') || '1');

// TODO:
const CONNECTION_POOL_SIZE = 10;

export function createConnectionPool(): ReturnType<typeof postgres> {
  return postgres({
    host: Deno.env.get('DATABASE_HOST_NAME'),
    database: Deno.env.get('DATABASE_NAME'),
    user: Deno.env.get('DATABASE_USERNAME'),
    password: Deno.env.get('DATABASE_PASSWORD'),
    port: parseInt(Deno.env.get('DATABASE_PORT') || '5432'),
    max: CONNECTION_POOL_SIZE,
  });
}
