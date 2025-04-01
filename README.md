# jvm-vs-jsr

Benchmarks comparing JVM and JS Runtime concurrency frameworks.

**CURRENTLY UNDER DEVELOPMENT**

## Purpose

To acquire a definitive understanding of how the throughput of I/O-bound work on Java virtual
threads compares with Java kernel threads, Java reactive frameworks, and JavaScript runtime
platforms (e.g. Node.js). Only PostgreSQL queries are examined.

## Introduction

When Node.js first came out, people were astonished that it could provide better throughput than
Java for I/O-bound operations. Numerous reactive frameworks have since been created to rectify the
issue (Java), but I had trouble finding apples-to-apples comparisons of how well these frameworks
perform relative to Node.js. TechEmpower [provides benchmarks](https://www.techempower.
com/benchmarks),
but I found it difficult to understand the commonalities and differences between any two
implementations. The present repo is an effort to get a definitive sense of how throughput compares
on these platforms for I/O-bound work, specifically PostgreSQL queries.

I'll be collecting throughput metrics for non-blocking APIs and for APIs that hit a backend
Postgres database. I'll be considering the
database clients JDBC (including HikariCP), R2DBC, and the Vert.x PG client.

I ended up writing an unexpected amount of Go client code for this project, but there are no
unit tests because I'm not planning to maintain the repo long term.

## The Plan

The first goal is to compare bare-metal Java solutions to bare-metal JavaScript runtime
solutions to learn the best possible throughput available. Fastify seems like a
reasonable bare-bones approach to testing on a JavaScript runtime, while Jooby seems like a
reasonable bare-bones approach for the JVM, given that Jooby code can be written in an
Express.js-like fashion. I also want to compare these results to Spring Boot, which seems like a
necessary baseline for comparing the performance of Java frameworks. I'll then look at how
Quarkus and Micronaut perform with R2DBC reactive I/O.

Here are the JavaScript Runtime combinations I'll be benchmarking:

- Fastify
  - Clustered Node.js + postgresjs
  - Clustered Deno + postgresjs
  - Clustered Bun + postgresjs

Here are the JVM combinations I'll be benchmarking:

- Jooby
  - JDBC with kernel threads (Java)
  - JDBC with virtual threads (Java)
  - R2DBC (Java)
  - Vert.x (Java)
- Spring Boot
  - JDBC with kernel threads (Java)
  - JDBC with virtual threads (Java)
  - WebFlux/R2DBC (Java)
  - WebFlux/R2DBC (Kotlin)
- Quarkus
  - R2DBC (Kotlin)
- Micronaut
  - R2DBC (Kotlin)

I'm benchmarking the reactive framework combinations in Kotlin because I find it unreasonable to
implement reactive frameworks for the JVM without coroutines, for the same reason it's
unreasonable to implement I/O in Node.js without async/await.

For each combination, I'll also be implemented a thread sleep to measure raw concurrency
unimpeded by waiting on a third tier.

## Installation and Setup

Install kubectl, helm, and helmfile, and configure kubectl for your cluster.

Run `helmfile init` to further install the Helm "diff" and "secrets" plugins.

Set the docker hostname image prefix in the following environment variable, excluding any trailing
+`/`:

```bash
export DOCKER_IMAGE_PREFIX=<your-docker-image-prefix>
```

And set it in `charts/values-secret.yaml` (which is in `.gitignore`):

```bash
global:
  app:
    dockerImagePrefix: <your-docker-image-prefix>
```

If your cluster does not automatically provision persistent volumes, you'll also need to create 
a volume called `client-postgres-volume`.

Note: I found cross-platform building using `buildx` too unreliable to use,
perhaps because of the dependency on calling `docker buildx create --use`.
The client image therefore builds only for `amd64`.

Label the three nodes as follows:

```bash
kubectl label nodes <node-1-name> kubernetes.io/hostname=client --overwrite
kubectl label nodes <node-2-name> kubernetes.io/hostname=app --overwrite
kubectl label nodes <node-3-name> kubernetes.io/hostname=backend --overwrite
```

Add the Helm repos for Prometheus and Grafana:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

## Building and Deploying

Create your cluster and configure `kubectl` to use it. Then:

```bash
mvn clean install
./bin/deploy backend-database
./bin/deploy client
./bin/deploy spring-jdbc-kernel-app # or another app
```

The `deploy` command deploys or redeploys. In the case of apps, it replaces the currently
deployed app (if any) with the named app.

## Running Benchmarks

**TODO: Obsolete/rewrite**

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Run `./benchmark setup <scenario>` to set up the scenario of the given name.
3. Run `./benchmark run <scenario> -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to see other commands and get usage help.

When running a test scenario, it outputs the first response for each unique combination
of status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body. This
output assists with debugging newly added applications.

## Useful Commands

```bash
./bin/undeploy backend-database
./bin/undeploy client
./bin/undeploy spring-jdbc-kernel-app # or another app
```
