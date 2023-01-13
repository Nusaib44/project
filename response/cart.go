package response

type CartProduct struct {
	ProductName   string
	Category      int
	Brand         string
	Price         int
	ProductOffer  int
	CategoryOffer int
	Total         int
	Image         string
	Description   string
}
type Cart struct {
	Total int
	Walet int
	// Product interface{}
}
