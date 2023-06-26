package account

type Account struct {
	Password string `json:"-"`
	UUID     string `json:"id"`
}

type AccountUpdate struct {
	Password string `json:"password"`
	Verify   string `json:"verify"`
}

func (a AccountUpdate) Patch() error {
	// todo verify password
	// todo update password
	return nil
}
