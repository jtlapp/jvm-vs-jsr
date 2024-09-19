import { Setup } from './order-items/setup.ts';

const DATABASE_URL = 'postgres://pgbouncer-service:6432/testdb';
const USERNAME = 'user';
const PASSWORD = 'password';

const setup = new Setup(DATABASE_URL, USERNAME, PASSWORD);
await setup.run();
console.log('Database setup completed.');