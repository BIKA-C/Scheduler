package account

type Academy struct {
	Account `json:"-"`
	Contact
	Name string `json:"name"`
}
