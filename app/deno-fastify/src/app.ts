import { installEndpoints } from "./endpoints.ts";
import { createConnectionPool, NUM_WORKERS, SERVER_PORT } from "./config.ts";

class App {
  async run() {
    // Deno doesn't use the same cluster module as Node.js
    // We'll use Deno.Worker instead for multi-threading
    if (Deno.env.get("DENO_WORKER") !== "true") {
      await this.startPrimaryThread();
    } else {
      await this.startWorkerThread();
    }
  }

  async startPrimaryThread() {
    console.log(`Primary ${Deno.pid} is running`);

    for (let i = 0; i < NUM_WORKERS; i++) {
      // Start worker processes
      const worker = new Deno.Command(Deno.execPath(), {
        args: ["run", "--allow-net", "--allow-env", "--allow-sys", "src/app.ts"],
        env: {
          DENO_WORKER: "true",
        },
        stdout: "inherit",
        stderr: "inherit",
      });
      
      const process = worker.spawn();
      
      // Handle process termination and restart
      process.status.then((status) => {
        console.log(`Worker process exited with code: ${status.code}`);
        // Restart worker
        this.startPrimaryThread();
      });
    }

    // Keep the primary process running
    await new Promise(() => {});
  }

  async startWorkerThread() {
    const pool = createConnectionPool();
    
    // Import Fastify for Deno
    const Fastify = await import("npm:fastify@5.1.0");
    const fastify = Fastify.default();
    
    // Install the endpoints
    installEndpoints(pool, fastify);

    try {
      // Start the server
      await fastify.listen({ port: SERVER_PORT, host: "0.0.0.0" });
      console.log(`Worker ${Deno.pid} started and listening on port ${SERVER_PORT}`);
    } catch (err) {
      console.error(err);
      Deno.exit(1);
    }
  }
}

// Run the application
await new App().run();
