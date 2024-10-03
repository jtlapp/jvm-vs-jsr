import postgres from 'postgresjs';

export class DatabaseUtils {
  static async dropTables(sql: ReturnType<typeof postgres>) {
    const result = await sql`
      SELECT tablename FROM pg_tables WHERE schemaname = 'public';
    `;
    for (const row of result) {
      if (row.tablename != 'shared_queries') {
        await sql`DROP TABLE IF EXISTS ${sql(row.tablename)} CASCADE`;
      }
    }
  }

  static async emptyTable(sql: ReturnType<typeof postgres>, tableName: string) {
    await sql`DELETE FROM ${sql(tableName)}`;
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
