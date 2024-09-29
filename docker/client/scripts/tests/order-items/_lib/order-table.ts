import { postgres } from '../../../_lib/deps.ts';

export class OrderTable {
  static async createTable(sql: ReturnType<typeof postgres>) {
    await sql`
        CREATE TABLE IF NOT EXISTS orders (
          id VARCHAR PRIMARY KEY,
          user_id VARCHAR REFERENCES users(id),
          order_date TIMESTAMP,
          status VARCHAR
        )
      `;
  }

  static createID(userID: string, orderNumber: number): string {
    return `${userID}_ORDER_${orderNumber}`;
  }

  static async insertOrder(
    sql: ReturnType<typeof postgres>,
    orderID: string,
    userID: string,
    orderDate: Date,
    status: string
  ) {
    await sql`
        INSERT INTO orders (id, user_id, order_date, status)
          VALUES (${orderID}, ${userID}, ${orderDate}, ${status})
      `;
  }
}
