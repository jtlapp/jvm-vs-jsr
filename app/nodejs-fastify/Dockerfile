FROM node:23-alpine

WORKDIR /app

COPY . .
RUN npm install

EXPOSE ${APP_PORT}
CMD ["node", "src/app.js"]