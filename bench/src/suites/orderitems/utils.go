package orderitems

import "fmt"

func toUserID(number int) string {
	return fmt.Sprintf("USER_%d", number)
}

func toProductID(number int) string {
	return fmt.Sprintf("PRODUCT_%d", number)
}

func toOrderID(userID string, orderNumber int) string {
	return fmt.Sprintf("%s_ORDER_%d", userID, orderNumber)
}

func toOrderItemID(orderID string, itemNumber int) string {
	return fmt.Sprintf("%s_ITEM_%d", orderID, itemNumber)
}
