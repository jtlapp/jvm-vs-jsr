import pkg from 'pg';
const { Pool } = pkg;

export const APP_NAME = process.env.APP_NAME;
export const APP_VERSION = '0.1.0';
export const SERVER_PORT = process.env.SERVER_PORT || 3000;
export const NUM_WORKERS = parseInt(process.env.NUM_WORKERS) || 1;

export function createConnectionPool() {
  return new Pool({
    host: process.env.DATABASE_HOST_NAME,
    database: process.env.DATABASE_NAME,
    user: process.env.DATABASE_USERNAME,
    password: process.env.DATABASE_PASSWORD,
    port: process.env.DATABASE_PORT,
  });
}
