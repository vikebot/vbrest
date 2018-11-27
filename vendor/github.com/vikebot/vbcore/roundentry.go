package vbcore

type Roundentry struct {
	Round
	Authtoken  string `json:"authtoken"`
	Watchtoken string `json:"watchtoken"`
}

type RoundentryConnectinfo struct {
	Roundticket string `json:"ticket"`
	AESKey      string `json:"aes_key"`
	IPv4        string `json:"ipv4"`
	IPv6        string `json:"ipv6"`
	Port        int    `json:"port"`
}

type RoundentryVerification struct {
	UserID  int     `json:"user_id"`
	RoundID int     `json:"round_id"`
	AESKey  *string `json:"aes_key,omitempty"`
}
