FROM denoland/deno:ubuntu AS build
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    unzip
RUN cd /tmp && \
    git clone https://github.com/wg/wrk.git
RUN cd /tmp/wrk && \
    make

FROM denoland/deno:ubuntu
COPY --from=build /tmp/wrk/wrk /usr/local/bin/
RUN apt-get update && apt-get install -y \
    libgcc-s1 \
    curl \
    vim
COPY root/ /root/
COPY scripts/ /scripts/
RUN find /scripts -type f -name ".*" -exec rm {} \;
CMD ["/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"]
WORKDIR /scripts
