package cart

import (
	"fmt"

	"github.com/zinx110/golang-backend-rest/types"
)

func getCartItemsIDs(items []types.CartItem) []int {
	ids := make([]int, 0, len(items))
	for _, item := range items {
		if item.ProductID == 0 {
			continue
		}
		ids = append(ids, item.ProductID)
	}
	return ids
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// check if all products are in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	// calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of product from db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		_, err := h.productStore.UpdateProduct(&product)
		if err != nil {
			return 0, 0, fmt.Errorf("error updating product %d: %w", product.ID, err)
		}

	}
	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "abc address", // In real scenario, get this from user profile or input

	})

	if err != nil {
		return 0, 0, fmt.Errorf("error creating order: %w", err)
	}

	// create order items
	for _, item := range items {
		product := productMap[item.ProductID]
		err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("error creating order item for product %d: %w", item.ProductID, err)
		}
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(items []types.CartItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product with id %d not found", item.ProductID)
		}
		if item.Quantity > product.Quantity {
			return fmt.Errorf("product %s is out of stock", product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	total := 0.0
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			continue
		}
		total += float64(item.Quantity) * product.Price
	}
	return total
}
