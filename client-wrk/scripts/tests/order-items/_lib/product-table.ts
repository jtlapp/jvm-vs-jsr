import postgres from 'postgresjs';
import { TestUtils } from 'test-utils';

export class ProductTable {
  static async createTable(sql: ReturnType<typeof postgres>) {
    await sql`
        CREATE TABLE IF NOT EXISTS products (
          id VARCHAR PRIMARY KEY,
          name VARCHAR,
          description TEXT,
          price NUMERIC,
          stock_quantity INTEGER,
          created_at TIMESTAMP DEFAULT NOW()
        )
      `;
  }

  static createID(productNumber: number): string {
    return TestUtils.createPaddedID('PRODUCT_', productNumber);
  }

  static async insertProduct(
    sql: ReturnType<typeof postgres>,
    productID: string,
    name: string,
    description: string,
    price: number,
    stockQuantity: number
  ) {
    await sql`
        INSERT INTO products (id, name, description, price, stock_quantity, created_at)
          VALUES (${productID}, ${name}, ${description}, ${price}, ${stockQuantity}, NOW())
      `;
  }
}
