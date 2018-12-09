package main

type conf struct {
	Addr string `json:"addr"`
	TLS  struct {
		Active bool   `json:"active"`
		Cert   string `json:"cert"`
		Key    string `json:"key"`
	} `json:"tls"`
	DB struct {
		Addr string `json:"addr"`
		User string `json:"user"`
		Pass string `json:"pass"`
		Name string `json:"name"`
	} `json:"db"`
	CORS struct {
		Wildcard bool `json:"wildcard"`
	} `json:"cors"`
	JWT struct {
		ProductionIsssuer   bool              `json:"production_issuer"`
		DefaultSigningKeyID string            `json:"default_signing_key_id"`
		SigningKeys         map[string]string `json:"signing_keys"`
	} `json:"jwt"`
	Recaptcha struct {
		Secret string `json:"secret"`
	} `json:"recaptcha"`
	Sendgrid struct {
		Secret string `json:"secret"`
	} `json:"sendgrid"`
}
