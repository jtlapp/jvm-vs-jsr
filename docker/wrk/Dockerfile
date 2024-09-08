FROM --platform=$BUILDPLATFORM alpine:latest as build
RUN apk add --update alpine-sdk perl linux-headers
RUN cd /tmp                               && \
    git clone https://github.com/wg/wrk.git
RUN cd /tmp/wrk                           && \
    make

FROM --platform=$TARGETPLATFORM alpine:latest
COPY --from=build /tmp/wrk/wrk /usr/local/bin/
RUN apk add --no-cache libgcc curl
WORKDIR /data
CMD ["/bin/ash", "-c", "exec /bin/ash -c 'trap : TERM INT; sleep 9999999999d & wait'"]