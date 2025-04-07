import { APP_NAME, APP_VERSION, NUM_WORKERS } from './config.ts';
import postgres from 'https://deno.land/x/postgresjs@v3.4.5/mod.js';

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export function installEndpoints(
  sql: ReturnType<typeof postgres>,
  server: any
) {
  server.get('/', async (request: any, reply: any) => {
    return 'Deno + Fastify + postgres';
  });

  server.get('/api/info', async (request: any, reply: any) => {
    return {
      appName: APP_NAME,
      appVersion: APP_VERSION,
      appConfig: {
        numWorkers: NUM_WORKERS,
      },
    };
  });

  server.get('/api/app-sleep', async (request: any, reply: any) => {
    const millis = parseInt(request.query.millis || '0');
    await delay(millis);

    return '{}';
  });

  server.get('/api/pg-sleep', async (request: any, reply: any) => {
    const millis = parseInt(request.query.millis || '0');
    const seconds = millis / 1000;

    await sql`SELECT pg_sleep(${seconds})`;
    return '{}';
  });
}
