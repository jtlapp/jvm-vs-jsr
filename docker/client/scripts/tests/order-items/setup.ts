import { AbstractSetup } from '../_lib/abstract-setup.ts';
import { SharedQueryRepo } from '../_lib/shared-query-repo.ts';

import { OrderItemTable } from './_lib/order-item-table.ts';
import { OrderTable } from './_lib/order-table.ts';
import { ProductTable } from './_lib/product-table.ts';
import { UserTable } from './_lib/user-table.ts';

const USER_COUNT = 1000;
const PRODUCT_COUNT = 700;
const ORDERS_PER_USER = 3;
const ITEMS_PER_ORDER = 4;

export class Setup extends AbstractSetup {
  constructor(dbURL: string, username: string, password: string) {
    super('order-items', dbURL, username, password);
  }

  protected async createTables() {
    await UserTable.createTable(this.sql);
    await ProductTable.createTable(this.sql);
    await OrderTable.createTable(this.sql);
    await OrderItemTable.createTable(this.sql);
  }

  protected async populateDatabase() {
    for (let i = 1; i <= USER_COUNT; i++) {
      await UserTable.insertUser(
        this.sql,
        UserTable.createID(i),
        `user${i}`,
        `user${i}@example.com`
      );
    }

    for (let i = 1; i <= PRODUCT_COUNT; i++) {
      await ProductTable.insertProduct(
        this.sql,
        ProductTable.createID(i),
        `Product ${i}`,
        `Description of product ${i}`,
        parseFloat((i % 50) + '.99'),
        100
      );
    }

    let orderedItemCount = 0;
    for (let i = 1; i <= USER_COUNT; i++) {
      for (let j = 1; j <= ORDERS_PER_USER; j++) {
        const userID = UserTable.createID(i);
        const orderID = OrderTable.createID(userID, j);
        await OrderTable.insertOrder(
          this.sql,
          orderID,
          userID,
          new Date(),
          'Shipped'
        );

        for (let k = 1; k <= ITEMS_PER_ORDER; k++) {
          const orderItemID = `${orderID}_ITEM_${k}`;
          const productNumber = (orderedItemCount % PRODUCT_COUNT) + 1;
          const productID = ProductTable.createID(productNumber);
          await OrderItemTable.insertOrderItem(
            this.sql,
            orderItemID,
            orderID,
            productID,
            1
          );
          orderedItemCount++;
        }
      }
    }
  }

  protected async createSharedQueries() {
    await SharedQueryRepo.createQuery(this.sql, {
      name: 'orderitems_getOrder',
      query: `
        SELECT o.id AS order_id, o.order_date, o.status, u.username, u.email,
               p.name, p.description, p.price, oi.quantity
        FROM orders o
        JOIN users u ON o.user_id = u.id
        JOIN order_items oi ON oi.order_id = o.id
        JOIN products p ON oi.product_id = p.id
        WHERE o.id = \${orderID}
      `,
      returns: 'rows',
    });

    await SharedQueryRepo.createQuery(this.sql, {
      name: 'orderitems_boostOrderItems',
      query: `
        UPDATE order_items oi
        SET quantity = quantity + 1
        FROM orders o
        WHERE oi.order_id = o.id AND o.id = \${orderID}
      `,
      returns: 'count',
    });
  }
}
