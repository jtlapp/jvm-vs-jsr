FROM alpine:latest as build
RUN apk add --update alpine-sdk perl linux-headers
RUN cd /tmp                               && \
    git clone -b 4.2.0 https://github.com/wg/wrk
RUN cd /tmp/wrk                           && \
    make

FROM alpine:latest
COPY --from=build /tmp/wrk/wrk /usr/local/bin/
RUN apk add --no-cache libgcc
WORKDIR /data
CMD ["/bin/ash", "-c", "exec /bin/ash -c 'trap : TERM INT; sleep 9999999999d & wait'"]