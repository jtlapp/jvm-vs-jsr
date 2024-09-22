import { postgres } from '../lib/deps.ts';

import { Utils } from '../lib/utils.ts';

export abstract class AbstractSetup {
  protected setupName: string;
  protected sql: ReturnType<typeof postgres>;

  constructor(
    setupName: string,
    dbURL: string,
    username: string,
    password: string
  ) {
    this.setupName = setupName;
    this.sql = postgres(dbURL, {
      username: username,
      password: password,
    });
  }

  getName(): string {
    return this.setupName;
  }

  async run(): Promise<void> {
    await Utils.dropTables(this.sql);
    await this.createTables();
    await this.populateDatabase();
    await this.createSharedQueries();
  }

  async recreateSharedQueries(): Promise<void> {
    await Utils.emptyTable(this.sql, 'shared_queries');
    await this.createSharedQueries();
  }

  async release(): Promise<void> {
    await this.sql.end();
  }

  protected abstract createTables(): Promise<void>;

  protected abstract populateDatabase(): Promise<void>;

  protected abstract createSharedQueries(): Promise<void>;
}
