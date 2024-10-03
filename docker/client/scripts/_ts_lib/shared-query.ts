import postgres from 'postgresjs';

/**
 * Properties of a query that is shared across the various benchmarked
 * API servers by being stored in a Postgres table.
 */
export interface SharedQueryProps {
  /**
   * Unique name by which servers can retrieve the query
   */
  name: string;
  /**
   * SQL query optionally containing parameters of the form `${paramName}`.
   */
  query: string;
  /**
   * What the query returns: 'nothing' (an empty response), 'rows' (a JSON
   * array), or 'count' (an integer row count).
   */
  returns: string;
}

/**
 * Query that is shared across the various benchmarked API servers.
 */
export class SharedQuery {
  private name: string; // unique name of the query in the database
  private query: string; // SQL query optionally with positional parameters
  private paramNames: string[] = []; // parameter names by positional index
  private returns: string; // what the query returns

  constructor(queryProps: SharedQueryProps) {
    this.name = queryProps.name;

    this.query = queryProps.query.replace(
      /\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}/g,
      (_, paramName) => {
        const paramIndex = this.paramNames.indexOf(paramName);
        if (paramIndex === -1) {
          this.paramNames.push(paramName);
          return `$${this.paramNames.length}`;
        }
        return `$${paramIndex + 1}`;
      }
    );

    this.returns = queryProps.returns;
  }

  /**
   * Returns the name of the shared query.
   */
  getName() {
    return this.name;
  }

  /**
   * Executes the query and returns its results, if any.
   * @param sql Postgres client for issuing queries
   * @param args An object mapping parameter names to their values.
   *      Must provide a defined/non-null argument for each parameter.
   * @return Nothing, an array of result rows, or a row count, depending
   *      on the value of the query's 'returns' property.
   */
  async execute(
    sql: ReturnType<typeof postgres>,
    args: { [key: string]: string | number | boolean }
  ) {
    const result = await sql.unsafe(
      this.query,
      this.paramNames.map((name) => args[name])
    );

    switch (this.returns) {
      case 'nothing':
        return;
      case 'rows':
        return result;
      case 'count':
        return result.count;
      default:
        throw new SharedQueryError(
          `Unrecognized returns type '${this.returns}'`
        );
    }
  }
}

export class SharedQueryError extends Error {
  constructor(message: string) {
    super(message);
  }
}
