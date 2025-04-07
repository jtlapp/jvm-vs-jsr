import { installEndpoints } from "./endpoints.ts";
import { createConnectionPool, NUM_WORKERS, SERVER_PORT } from "./config.ts";

class App {
  async run() {
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

    await new Promise(() => {});
  }

  async startWorkerThread() {
    const Fastify = await import("npm:fastify@5.1.0");
    const fastify = Fastify.default();
    
    const sql = createConnectionPool();
    installEndpoints(sql, fastify);

    try {
      await fastify.listen({ port: SERVER_PORT, host: "0.0.0.0" });
      console.log(`Worker ${Deno.pid} started and listening on port ${SERVER_PORT}`);
    } catch (err) {
      console.error(err);
      Deno.exit(1);
    }
  }
}

await new App().run();
