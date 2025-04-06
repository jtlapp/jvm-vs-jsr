import { APP_NAME, APP_VERSION, NUM_WORKERS } from "./config.ts";
import { Pool } from "https://deno.land/x/postgres@v0.17.0/mod.ts";

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export function installEndpoints(pool: Pool, server: any) {
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

    const client = await pool.connect();
    try {
      await client.queryObject('SELECT pg_sleep($1)', [seconds]);
      return '{}';
    } finally {
      client.release();
    }
  });
}
