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

TODO: Look again into using environment variables in the POM.

## Building and Deploying

TODO: Note about creating persistent docker volume.

Deploy the backend database and client for benchmarking the variety of apps. You can only 
install one app at a time.

```bash
mvn clean install
./bin/setup-network # creates the docker network with backend
cd bench/src
./benchmark setup-results # creates the benchmark results database
```

```bash
./bin/deploy spring-jdbc-kernel-app # or another app
```

The `deploy` command deploys or redeploys. In the case of apps, it replaces the currently
deployed app (if any) with the named app. You can undeploy as follows:

```bash
./bin/teardown app # undeploy just the app
./bin/teardown all # remove all container and the docker network, but not the results volume
```

## Running Benchmarks

1. Exec into the client pod using bash: `kubectl exec -it <client-pod> -- bash`.
2. Run `./benchmark setup-results` to create the benchmark results database.
3. Run `./benchmark setup-backend -scenario <scenario>` to set up the scenario of the given name.
   Only required if the scenario uses backend Postgres tables.
4. 
5. TODO: introduce try, run, and loop
6. Run `./benchmark run -scenario <scenario> -rate <requests-per-sec> -duration <seconds>`.

Run `./benchmark` to see other commands and get usage help.

When running a test scenario, it outputs the first response for each unique combination
of status code, shared query name, and error message. For queries erroneously returning
non-JSON, it also prints each unique combination of status code and response body. This
output assists with debugging newly added applications.
