package account

type Institution struct {
	Account `json:"account"`
	Contact Contact `json:"contact"`
	Name    string  `json:"name"`
}
