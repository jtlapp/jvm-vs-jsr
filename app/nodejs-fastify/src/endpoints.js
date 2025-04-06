import { APP_NAME, APP_VERSION, NUM_WORKERS } from './config.js';

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export function installEndpoints(pool, server) {
  server.get('/', async (request, reply) => {
    return 'Node.js + Fastify + pg package';
  });

  server.get('/api/info', async (request, reply) => {
    return {
      appName: APP_NAME,
      appVersion: APP_VERSION,
      appConfig: {
        numWorkers: NUM_WORKERS,
      },
    };
  });

  server.get('/api/app-sleep', async (request, reply) => {
    const millis = parseInt(request.query.millis || '0');
    await delay(millis);

    return '{}';
  });

  server.get('/api/pg-sleep', async (request, reply) => {
    const millis = parseInt(request.query.millis || '0');
    const seconds = millis / 1000;

    await pool.query('SELECT pg_sleep($1)', [seconds]);
    return '{}';
  });
}
