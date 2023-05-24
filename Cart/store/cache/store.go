package cache

import "github.com/faagerholm/coffee/cart/store"

type Cache struct {
	carts map[int]store.Cart
}

func New() *Cache {
	carts := make(map[int]store.Cart)
	return &Cache{
		carts: carts,
	}
}

func (c *Cache) AddItem(userID int, item store.Item) error {
	if cart, ok := c.carts[userID]; ok {
		cart.Items = append(cart.Items, item)
		c.carts[userID] = cart
	} else {
		c.carts[userID] = store.Cart{
			Items: []store.Item{item},
		}
	}
	return nil
}

func (c *Cache) GetCart(userID int) (store.Cart, error) {
	return c.carts[userID], nil
}

func (c *Cache) EmptyCart(userID int) error {
	delete(c.carts, userID)
	return nil
}
