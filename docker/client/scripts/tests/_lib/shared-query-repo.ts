import { postgres } from '../../_lib/deps.ts';

import {
  SharedQuery,
  SharedQueryProps,
  SharedQueryError,
} from './shared-query.ts';

/**
 * Repository of queries shared across the various benchmarked API
 * servers, backed by a postgres table.
 */
export class SharedQueryRepo {
  static _cachedQueries: { [key: string]: SharedQuery } = {};

  /**
   * Creates and returns a shared query having the given query properties,
   * storing it in the database under its provided name.
   * @param sql Postgres client for issuing queries
   * @param queryProps Properties of the shared query to create
   */
  static async createQuery(
    sql: ReturnType<typeof postgres>,
    queryProps: SharedQueryProps
  ) {
    if (this._cachedQueries[queryProps.name] !== undefined) {
      throw new Error(
        `Shared query of name '${queryProps.name}' already defined.`
      );
    }

    await sql`DELETE FROM shared_queries WHERE name = ${queryProps.name}`;

    const query = new SharedQuery(queryProps);
    await sql`
        INSERT INTO shared_queries (name, query, returns) 
          VALUES (${queryProps.name}, ${queryProps.query}, ${queryProps.returns})
      `;
    return query;
  }

  /**
   * Loads the shared query of the given name from the database and returns
   * it, returning a cached instance if it was already previously loaded.
   * @param sql Postgres client for issuing queries
   * @param queryName Name uniquely identifying query among all shared queries
   */
  static async loadQuery(
    sql: ReturnType<typeof postgres>,
    queryName: string
  ): Promise<SharedQuery> {
    let query = this._cachedQueries[queryName];
    if (query !== undefined) return query;

    const result = await sql<SharedQueryProps[]>`
      SELECT * FROM shared_queries WHERE name = ${queryName}
    `;
    if (result.length == 0) {
      throw new SharedQueryError(`SharedQuery named '${queryName}' not found.`);
    }
    query = new SharedQuery(result[0]);
    this._cachedQueries[queryName] = query;
    return query;
  }
}
