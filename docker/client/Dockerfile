# deno:latest sometimes does not work (e.g. ubuntu-2.0.0-rc.7)
FROM denoland/deno:ubuntu-2.0.0-rc.5 AS build
RUN apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    git \
    wget \
    unzip

# build with a clone of a repo that will build on modern linux
RUN cd /tmp && \
    git clone https://github.com/jtlapp/wrk2.git
RUN cd /tmp/wrk2 && \
    make

FROM denoland/deno:ubuntu
COPY --from=build /tmp/wrk2/wrk /usr/local/bin/
RUN apt-get update && apt-get install -y \
    libgcc-s1 \
    curl \
    vim
COPY root/ /root/
COPY scripts/ /scripts/
RUN find /scripts -type f -name ".*" -exec rm {} \;
CMD ["/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"]
WORKDIR /scripts
