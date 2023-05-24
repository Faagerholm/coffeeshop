package store

type Item struct {
	ID       int
	Name     string
	Price    float32
	Quantity int
}

type Cart struct {
	Items []Item
}

type Store interface {
	AddItem(userID int, item Item) error
	GetCart(userID int) (Cart, error)
	EmptyCart(userID int) error
}
