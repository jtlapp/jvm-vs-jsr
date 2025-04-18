import cluster from 'node:cluster';
import fastify from 'fastify';
import { NUM_WORKERS, SERVER_PORT, createConnectionPool } from './config.js';
import { installEndpoints } from './endpoints.js';

class App {
  async run() {
    if (NUM_WORKERS > 1 && cluster.isPrimary) {
      await this.startPrimaryThread();
    } else {
      await this.startWorkerThread();
    }
  }

  async startPrimaryThread() {
    console.log(`Primary ${process.pid} is running`);

    for (let i = 0; i < NUM_WORKERS; i++) {
      cluster.fork();
    }

    cluster.on('exit', (worker, code, signal) => {
      console.log(`Worker ${worker.process.pid} died`);
      cluster.fork();
    });
  }

  async startWorkerThread() {
    const sql = createConnectionPool();
    const server = fastify();
    installEndpoints(sql, server);

    try {
      server.listen({ port: SERVER_PORT, host: '0.0.0.0' });
      console.log(`Worker ${process.pid} started and listening`);
    } catch (err) {
      server.log.error(err);
      process.exit(1);
    }
  }
}

new App().run();
