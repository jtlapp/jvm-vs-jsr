# jvm-vs-jsr

Benchmarks comparing JVM and JS Runtime concurrency frameworks.

**CURRENTLY UNDER DEVELOPMENT**

## Introduction

When Node.js first came out, people were astonished that it could provide better throughput than 
Java for I/O-bound operations. Numerous reactive frameworks have since been created to rectify the 
issue in Java, but I had trouble finding apples-to-apples comparisons of how well these frameworks
perform relative to Node.js. TechEmpower [provides benchmarks](https://www.techempower.
com/benchmarks),
but I found it difficult to understand the commonalities and differences between any two
implementations. The present repo is my effort to determine which of the two platforms 
are better for I/O-bound work.

My objective is to select a high performance platform for REST API microservices. Go would be an 
obvious choice for the programming language, but I have not been happy with Go. I've narrowed the
language down to one of TypeScript, Java, and Kotlin.

I want to use a well-supported batteries-included web server framework. For Java and Kotlin, I 
plan to evaluate Spring Boot, Quarkus, and Micronaut. I'll be evaluating them with kernel 
threads, with virtual threads, and with reactive I/O. I also plan to evaluate the virtual-thread-only framework Helidon. I'll also be investigating GraalVM.

For TypeScript, I plan to evaluate Nest.js and probably also tsoa, running them on Node.js, Deno,
and Bun. I'll also be looking at improving throughput with worker threads.

I'll be collecting throughput metrics for non-blocking APIs and for APIs that hit a backend 
Postgres database. I'll be considering the 
database clients JDBC, R2DBC, and the Vert.x PG client.

I ended up writing an unexpected amount of Go client code for this project, but there are no 
unit tests because I don't plan to maintain the repo after making my decision.

## The Plan

There are too many combinations of the above configurations to test all of them, so I need to be 
strategic about proceeding. We can limit the use of Tomcat and Netty according to requirements
or recommendations. I'll benchmark each of the following scenarios:

- Spring Boot
  - Kernel threads (Tomcat) [APP]
    - non-blocking API
    - JDBC
  - Virtual threads (Tomcat) [APP]
    - non-blocking API
    - JDBC
  - Reactive IO (Netty)
    - R2DBC only (to establish need for WebFlux) [APP]
      - non-blocking API
      - postgres
    - WebFlux/R2DBC [APP]
      - non-blocking API
      - postgres
- Quarkus
  - Virtual threads (Netty) [APP]
    - non-blocking API
    - JDBC
  - Reactive IO (Netty)
    - non-blocking API
    - R2DBC [APP]
    - Vert.x PG [APP]
- Micronaut
  - Virtual threads (Netty) [APP]
    - non-blocking API
    - JDBC
  - Reactive IO (Netty)
    - non-blocking API
    - R2DBC [APP]
    - Vert.x PG [APP]
- Helidon MP (virtual threads, Nima web server) [APP]
    - non-blocking API
    - Helidon DB (non-blocking wrap of JDBC)
- Nest.js
  - Node.js
    - non-blocking API
    - pg (only works on Node.js) [APP]
    - postgresjs [APP]
  - Deno [APP]
    - non-blocking API
    - postgresjs
  - Bun [APP]
    - non-blocking API
    - postgresjs
- tsoa
  - Node.js
    - non-blocking API
    - pg (only works on Node.js) [APP]
    - postgresjs [APP]
  - Deno [APP]
    - non-blocking API
    - postgresjs
  - Bun [APP]
    - non-blocking API
    - postgresjs

The "[APP]" notation indicates that the combination represents a distinct application.

After establishing the most performant JVM and JS runtime scenarios, I'll explore improving them with GraalVM and worker threads, and I'll look at simplifying reactive code with Kotlin.

## Installation and Setup

Install kubectl, helm, and helmfile, and configure kubectl for your cluster.

Run `helmfile init` to further install the Helm "diff" and "secrets" plugins.

Set the docker image prefix in an environment variables:

```bash
export DOCKER_IMAGE_PREFIX=<your-docker-image-prefix>
```

And set it in `charts/values-secret.yaml` (which is in `.gitignore`):

```bash
global:
  app:
    dockerImagePrefix: <your-docker-image-prefix>
```

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
./bin/deploy database
./bin/deploy client
./bin/deploy spring-jdbc-kernel # or another app
```

The `deploy` command deploys or redeploys. In the case of apps, it replaces the currently 
deployed app (if any) with the named app.

## Running Benchmarks

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
./bin/undeploy database
./bin/undeploy client
./bin/undeploy spring-jdbc-kernel # or another app
```
