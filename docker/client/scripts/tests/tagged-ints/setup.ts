import { SharedQueryRepo } from '../_lib/shared-query-repo.ts';

import { TaggedIntTable } from './_lib/tagged-int-table.ts';
import { TestUtils } from '../_lib/test-utils.ts';
import { AbstractSetup } from '../_lib/abstract-setup.ts';

const ROW_COUNT = 1000000;
const MAX_INT = 1000;
const RANDOM_SEEDS = [12345, 54321, 98765, 56789];

export class Setup extends AbstractSetup {
  static TAG_CHARS = '0123456789ABCDEF';
  static TAG_CHARS_LENGTH = Setup.TAG_CHARS.length;

  constructor(dbURL: string, username: string, password: string) {
    super('tagged-ints', dbURL, username, password);
  }

  protected async createTables() {
    await TaggedIntTable.createTable(this.sql);
  }

  protected async populateDatabase() {
    const random = TestUtils.createRandomNumberGenerator(RANDOM_SEEDS);

    for (let i = 1; i <= ROW_COUNT; i++) {
      await TaggedIntTable.insertTaggedInt(
        this.sql,
        this.createTag(random),
        this.createTag(random),
        Math.floor(random() * MAX_INT)
      );
    }
  }

  protected async createSharedQueries() {
    await SharedQueryRepo.createQuery(this.sql, {
      name: 'taggedints_sumInts',
      query: `
        SELECT SUM(int) AS sum
        FROM tagged_ints
        WHERE tag1 = \${tag1} AND tag2 = \${tag2}
      `,
      returns: 'rows',
    });

    await SharedQueryRepo.createQuery(this.sql, {
      name: 'taggedints_getInt',
      query: `SELECT int FROM tagged_ints WHERE id = \${id}`,
      returns: 'rows',
    });
  }

  private createTag(random: () => number): string {
    return (
      Setup.TAG_CHARS[Math.floor(random() * Setup.TAG_CHARS_LENGTH)] +
      Setup.TAG_CHARS[Math.floor(random() * Setup.TAG_CHARS_LENGTH)]
    );
  }
}
