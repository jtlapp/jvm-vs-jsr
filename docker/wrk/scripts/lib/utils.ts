import { postgres } from './deps.ts';

export class Utils {
  private static ZERO_PADDING_WIDTH = 6;

  static createPaddedID(prefix: string, value: number): string {
    return prefix + value.toString().padStart(this.ZERO_PADDING_WIDTH, '0');
  }

  static async dropTables(sql: ReturnType<typeof postgres>) {
    const result = await sql`
      SELECT tablename FROM pg_tables WHERE schemaname = 'public';
    `;
    for (const row of result) {
      if (row.tableName != 'queries') {
        await sql`DROP TABLE IF EXISTS ${row.tableName} CASCADE`;
      }
    }
  }

  /**
   * Returns a postgres client for accessing the database.
   */
  static openClient() {
    return postgres({
      host: 'pgbouncer-service',
      port: 6432,
      database: 'testdb',
      username: 'user',
      password: 'password',
    });
  }
}
