# jvm-vs-jsr

Benchmarks comparing JVM and JS Runtime concurrency frameworks

**CURRENTLY UNDER DEVELOPMENT**

## Purpose

To acquire a definitive understanding of how the throughput of I/O-bound work on Java virtual
threads compares with Java kernel threads, Java reactive frameworks, and JavaScript runtime
platforms (e.g. Node.js). Only PostgreSQL queries are examined.

## Introduction

When Node.js first came out, people were astonished that it could provide better throughput than
Java for I/O-bound operations. Numerous reactive frameworks have since been created to rectify the
issue (Java), but I had trouble finding apples-to-apples comparisons of how well these frameworks
perform relative to Node.js and to each other. TechEmpower [provides benchmarks](
https://www.techempower.com/benchmarks), but the implementations are so dramatically different
that I didn't understand how to compare them. The present repo is an effort to get a definitive
sense of how throughput compares on these platforms for I/O-bound work, specifically PostgreSQL
queries.

The benchmark is performed all on a single machine via docker containers. The benchmark client 
hits an app container, which hits a pgBouncer container, which hits a Postgres container. I 
proxied Postgres with pgBouncer to allow Postgres to accept as many concurrent connections as the 
client can make, to reduce the effect of Postgres itself on app concurrency. I'll be benchmarking
the database clients JDBC (including HikariCP), R2DBC, and the Vert.x PG.

I ended up writing an unexpected amount of Go client code for this project, but there are no
unit tests because at present I'm not planning to maintain the repo long term.

## The Comparisons

The first goal is to compare bare-metal Java solutions to bare-metal JavaScript runtime
solutions to learn the best possible throughput available. Fastify seems like a
reasonable bare-bones approach to testing on a JavaScript runtime, while Jooby seems like a
reasonable bare-bones approach for the JVM, given that Jooby code can be written in an
Express.js-like fashion. I also want to compare these results to Spring Boot, which seems like a
necessary baseline for comparing the performance of Java frameworks. I'll then look at how
Quarkus and Micronaut perform with R2DBC reactive I/O.

Here are the JavaScript Runtime combinations I'll be benchmarking:

- Fastify
  - Clustered Node.js + postgresjs (JavaScript)
  - Clustered Deno + postgresjs (JavaScript)
  - Clustered Bun + postgresjs (JavaScript)

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

## Tooling Installation

- Go
- Java
- Kotlin

TODO: Convert the following to Ubunto installation instructions:
```
ENV GO_VERSION=1.23.2
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN apt-get update && \
    apt-get install -y \
        vim \
        curl \
        net-tools \
        iproute2 \
        git \
        build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN curl -OL https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

RUN mkdir -p $GOPATH/bin

COPY root/ /root/
COPY src/ /src/

# remove builds made on my development machine
RUN find /src -type f -name "benchmark" -exec rm {} \;

WORKDIR /src
RUN go build
```

Useful `~/.vimrc`:
```
set tabstop=2
set shiftwidth=2
set softtabstop=2
set expandtab
```

## Configuration

The configuration is largely provided by environment variables, with `config/env-config` 
setting them to default values. Run this script or a modified copy of the script before running 
any of the `bin/` scripts and before running the benchmark tool. 

If you decide to change the default docker image prefix, in addition to changing the 
`DOCKER_IMAGE_PREFIX` environment variable, you'll need to assign the prefix via
`-Ddocker-image-prefix=<prefix>` when you build the applications.

By default, the applications build for AMD64, but you can use `-Ddocker-architecture=arm64` to 
build for Apple Silicon.

You may need to tweak `BENCH_MAX_RESERVED_PORTS` for your machine. The benchmark makes heavy use 
of ports and so checks to make sure the ports are ready before each test. One of the checks 
verifies that the number of established ports (ports having active connections) is at most
`BENCH_MAX_RESERVED_PORTS`. You'll want to set this to as low a number as you can for maximum 
confidence that the machine is ready to perform a benchmark. To see the current number of 
established ports, set this variable to 0 and run `bench/src/benchmark`. The tool with terminate 
with an error telling you the number of ports currently in use.

TODO: Look again into using environment variables in the POM.

TODO: Either mention need to kill benchmark to abort, or add interrupt support.

## Building

To build everything, including the `benchmark` tool, run the following from the repo root:

```bash
mvn clean install
```

To only rebuild the `benchmark` tool:

```bash
cd bench/src
make
```

To only rebuild one app:

```bash
mvn clean install -pl app/<app-directory> -am
```

## Deploying

After building, to deploy the docker network of containers for use with benchmarking, run the 
following script:

```bash
./bin/setup-network
```

This creates containers for pgBouncer, the backend database, and the results database. 
`config/env-config` specifies the names of these containers. The results database uses a 
persistent volume whose name is also given in `env-config`. The results database keeps a record 
of all benchmark runs and its data persists until you delete the volume via docker.

Upon first creating the persistent volume, set up the results database as follows. You won't 
need to run this command again until you replace the persistent volume:

```bash
./bench/src/benchmark setup-results # creates the benchmark results database
```

To run a benchmark for an app, you have to first deploy that app. Only one app can be deployed 
at a time. The script for deploying an app undeploys any already-existing app before deploying 
the requested app:

```bash
./bin/deploy spring-jdbc-kernel-app # or another app
```

Run `./bin/deploy` with no arguments to see the available apps. These apps roughly correspond to 
the names of the directories within `app/`, but each of these directories defines a server that 
can be variously configured, so rely on `./bin/deploy` (no args) for the named configurations.

## Teardown

To remove the container for just the current app, whatever app is running, do the following 
(`app` here is a literal, not the name of the app):

```bash
./bin/teardown app
```

To remove all containers, including the app container, and the docker network too, run this:

```bash
./bin/teardown all
```

The above command does not remove the volume holding the results database. To remove this volume 
(discarding its data), you'll need to do so via docker directly.

## Designing Benchmarks

TODO: discuss common config params

## Running Benchmarks

TODO:

Run benchmarks via the `benchmark` tool, which is in `bench/src`. You'll likely want to `cd` to 
this directory for convenience.

3. Run `./benchmark setup-backend -scenario <scenario>` to set up the scenario of the given name.
   Only required if the scenario uses backend Postgres tables.
5. TODO: introduce try, run, and loop
6. Run `./benchmark run -scenario <scenario> -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to see other commands and get usage help.

When running a test scenario, it outputs the first response for each unique combination
of status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body. This
output assists with debugging newly added applications.
