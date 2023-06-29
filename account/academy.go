package account

type Academy struct {
	Account
	Contact Contact `json:"contact"`
	Name    string  `json:"name"`
}
