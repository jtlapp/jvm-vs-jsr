import postgres from 'postgresjs';

export class TaggedIntTable {
  static async createTable(sql: ReturnType<typeof postgres>) {
    await sql`
        CREATE TABLE IF NOT EXISTS tagged_ints (
          id BIGSERIAL PRIMARY KEY,
          tag1 VARCHAR NOT NULL,
          tag2 VARCHAR NOT NULL,
          int INTEGER NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )
      `;
  }

  static async insertTaggedInt(
    sql: ReturnType<typeof postgres>,
    tag1: string,
    tag2: string,
    int: number
  ) {
    await sql`
        INSERT INTO tagged_ints (tag1, tag2, int, created_at)
          VALUES (${tag1}, ${tag2}, ${int}, NOW())
      `;
  }
}
