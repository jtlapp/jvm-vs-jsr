const cluster = require('cluster');
const numCPUs = require('os').cpus().length;
const fastify = require('fastify');
const { Pool } = require('pg');

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

  for (let i = 0; i < numCPUs; i++) {
    cluster.fork();
  }

  cluster.on('exit', (worker, code, signal) => {
    console.log(`Worker ${worker.process.pid} died`);
    cluster.fork();
  });
} else {
  const server = fastify({
    logger: {
      transport: {
        target: 'pino-pretty',
      },
      level: 'error',
    },
  });

  server.get('/', async (request, reply) => {
    return 'Node.js + Fastify + pg package';
  });

  server.get('/api/app-info', async (request, reply) => {
    return {
      appName: process.env.APP_NAME,
      appVersion: process.env.APP_VERSION,
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
      await server.listen({ port: 3000, host: '0.0.0.0' });
      console.log(`Worker ${process.pid} started and listening`);
    } catch (err) {
      server.log.error(err);
      process.exit(1);
    }
  };

  start();
}
