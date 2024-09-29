import { postgres } from '../../../_lib/deps.ts';

export class OrderItemTable {
  static async createTable(sql: ReturnType<typeof postgres>) {
    await sql`
        CREATE TABLE IF NOT EXISTS order_items (
          id VARCHAR PRIMARY KEY,
          order_id VARCHAR REFERENCES orders(id),
          product_id VARCHAR REFERENCES products(id),
          quantity INTEGER
        )
      `;
  }

  static async insertOrderItem(
    sql: ReturnType<typeof postgres>,
    orderItemID: string,
    orderID: string,
    productID: string,
    quantity: number
  ) {
    await sql`
        INSERT INTO order_items (id, order_id, product_id, quantity)
          VALUES (${orderItemID}, ${orderID}, ${productID}, ${quantity})
      `;
  }
}
