import { postgres } from '../../../_lib/deps.ts';
import { TestUtils } from '../../_lib/test-utils.ts';

export class UserTable {
  static async createTable(sql: ReturnType<typeof postgres>) {
    await sql`
        CREATE TABLE IF NOT EXISTS users (
          id VARCHAR PRIMARY KEY,
          username VARCHAR UNIQUE NOT NULL,
          email VARCHAR UNIQUE NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )
      `;
  }

  static createID(userNumber: number): string {
    return TestUtils.createPaddedID('USER_', userNumber);
  }

  static async insertUser(
    sql: ReturnType<typeof postgres>,
    userID: string,
    username: string,
    email: string
  ) {
    await sql`
        INSERT INTO users (id, username, email, created_at)
          VALUES (${userID}, ${username}, ${email}, NOW())
      `;
  }
}
