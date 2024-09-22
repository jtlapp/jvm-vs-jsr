import { Setup } from './tagged-ints/setup.ts';

const DATABASE_URL = 'postgres://pgbouncer-service:6432/testdb';
const USERNAME = 'user';
const PASSWORD = 'password';

const setup = new Setup(DATABASE_URL, USERNAME, PASSWORD);
await setup.run();
// await setup.recreateSharedQueries();
await setup.release();
console.log(`'${setup.getName()}' database setup completed.`);
