const cluster = require('cluster');
const fastify = require('fastify');
const { Pool } = require('pg');

const APP_NAME = process.env.APP_NAME;
const APP_VERSION = '0.1.0';
const SERVER_PORT = process.env.SERVER_PORT;

const NUM_WORKERS = parseInt(process.env.NUM_WORKERS) || 1;

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const pool = new Pool({
  host: process.env.DATABASE_HOST_NAME,
  database: process.env.DATABASE_NAME,
  user: process.env.DATABASE_USER,
  password: process.env.DATABASE_PASSWORD,
  port: process.env.DATABASE_PORT,
});

if (cluster.isPrimary) {
  console.log(`Primary ${process.pid} is running`);

  for (let i = 0; i < NUM_WORKERS; i++) {
    cluster.fork();
  }

  cluster.on('exit', (worker, code, signal) => {
    console.log(`Worker ${worker.process.pid} died`);
    cluster.fork();
  });
} else {
  const server = fastify({});

  server.get('/', async (request, reply) => {
    return 'Node.js + Fastify + pg package';
  });

  server.get("/api/info", async (request, reply) => {
    return {
      appName: APP_NAME,
      appVersion: APP_VERSION,
      appConfig: {
        numWorkers: NUM_WORKERS
      }
    };
  });

  server.get('/api/app-info', async (request, reply) => {
    return {
      appName: APP_NAME,
      appVersion: APP_VERSION,
      appConfig: {},
    };
  });

  server.get('/api/app-sleep', async (request, reply) => {
    const millis = parseInt(request.query.millis);
    await delay(millis);

    return '{}';
  });

  server.get('/api/pg-sleep', async (request, reply) => {
    const millis = parseInt(request.query.millis);
    const seconds = millis / 1000;

    await pool.query('SELECT pg_sleep($1)', [seconds]);
    return '{}';
  });

  const start = async () => {
    try {
      await server.listen({ port: SERVER_PORT, host: '0.0.0.0' });
      console.log(`Worker ${process.pid} started and listening`);
    } catch (err) {
      server.log.error(err);
      process.exit(1);
    }
  };

  start();
}
