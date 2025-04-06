FROM denoland/deno:alpine-2.2.8

WORKDIR /app

COPY . .

RUN deno cache --unstable-worker-options \
    --node-modules-dir \
    src/app.ts

EXPOSE ${APP_PORT}

CMD ["deno", "run", \
     "--allow-net", \
     "--allow-env", \
     "--allow-run", \
     "--allow-read", \
     "--allow-write", \
     "--unstable-worker-options", \
     "--node-modules-dir", \
     "src/app.ts"]
