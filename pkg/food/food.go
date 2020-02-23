package food

// Food is one food item at a canteen.
type Food struct {
	Name       string `json:"name"`
	StudPrice  int `json:"studprice"`
	ProfPrice  int `json:"profprice"`
	Vegan      bool `json:"vgn"`
	Vegetarian bool `json:"vgt"`
}
