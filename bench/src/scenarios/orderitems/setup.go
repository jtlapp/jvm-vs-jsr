package orderitems

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	totalUsers    = 1000
	totalProducts = 700
	ordersPerUser = 3
	itemsPerOrder = 4
)

type SetupImpl struct {
	dbPool *pgxpool.Pool
}

func (s *SetupImpl) CreateTables() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
          id VARCHAR PRIMARY KEY,
          username VARCHAR UNIQUE NOT NULL,
          email VARCHAR UNIQUE NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )`
	_, err := s.dbPool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	query = `
        CREATE TABLE IF NOT EXISTS products (
          id VARCHAR PRIMARY KEY,
          name VARCHAR,
          description TEXT,
          price NUMERIC,
          stock_quantity INTEGER,
          created_at TIMESTAMP DEFAULT NOW()
        )`
	_, err = s.dbPool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	query = `
        CREATE TABLE IF NOT EXISTS orders (
          id VARCHAR PRIMARY KEY,
          user_id VARCHAR REFERENCES users(id),
          order_date TIMESTAMP,
          status VARCHAR
        )`
	_, err = s.dbPool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	query = `
        CREATE TABLE IF NOT EXISTS order_items (
          id VARCHAR PRIMARY KEY,
          order_id VARCHAR REFERENCES orders(id),
          product_id VARCHAR REFERENCES products(id),
          quantity INTEGER
        )`
	_, err = s.dbPool.Exec(context.Background(), query)
	return err
}

func (s *SetupImpl) PopulateTables() error {
	for i := 1; i <= totalUsers; i++ {
		query := `INSERT INTO users (id, username, email, created_at) VALUES ($1, $2, $3, NOW())`
		username := fmt.Sprintf("user%d", i)
		_, err := s.dbPool.Exec(context.Background(), query, toUserID(i), username, username+"@example.com")
		if err != nil {
			return err
		}
	}

	for i := 1; i <= totalProducts; i++ {
		query := `INSERT INTO products (id, name, description, price, stock_quantity, created_at)
					VALUES ($1, $2, $3, $4, $5, NOW())`
		productName := fmt.Sprintf("product-%d", i)
		_, err := s.dbPool.Exec(context.Background(), query, toProductID(i), productName,
			productName+" description", float64(i%50)+0.99, 100)
		if err != nil {
			return err
		}
	}

	orderedItemCount := 0
	for i := 1; i <= totalUsers; i++ {
		userID := toUserID(i)

		for j := 1; j <= ordersPerUser; j++ {
			orderID := toOrderID(userID, j)

			query := `INSERT INTO orders (id, user_id, order_date, status)
						VALUES ($1, $2, NOW(), 'shipped')`
			_, err := s.dbPool.Exec(context.Background(), query, toOrderID(userID, j), userID)
			if err != nil {
				return err
			}

			for k := 1; k <= itemsPerOrder; k++ {
				query := `INSERT INTO order_items (id, order_id, product_id, quantity)
							VALUES ($1, $2, $3, $4)`
				orderItemID := toOrderItemID(orderID, k)
				productNumber := (orderedItemCount % totalProducts) + 1
				productID := toProductID(productNumber)
				_, err := s.dbPool.Exec(context.Background(), query, orderItemID, orderID, productID, 1)
				if err != nil {
					return err
				}
				orderedItemCount++
			}
		}
	}
	return nil
}

func (s *SetupImpl) GetSharedQueries() []util.SharedQuery {
	return []util.SharedQuery{
		{
			Name: "orderitems_getOrder",
			Query: `
				SELECT o.id AS order_id, o.order_date, o.status, u.username, u.email,
						p.name, p.description, p.price, oi.quantity
					FROM orders o
					JOIN users u ON o.user_id = u.id
					JOIN order_items oi ON oi.order_id = o.id
					JOIN products p ON oi.product_id = p.id
					WHERE o.id = ${orderID}
			`,
			Returns: "rows",
		},
		{
			Name: "orderitems_boostOrderItems",
			Query: `
				UPDATE order_items oi
					SET quantity = quantity + 1
					FROM orders o
					WHERE oi.order_id = o.id AND o.id = ${orderID}
			`,
			Returns: "count",
		},
	}
}
