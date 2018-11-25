package configure

// AppConfig :
type AppConfig struct {
	StaticFilePath    string
	LoginRedirectURL  string
	LogoutRedirectURL string
	CookieSecretKey   string
}

// AppProperties :
var AppProperties AppConfig

// SetProperties :
func SetProperties(ap AppConfig) {
	AppProperties = ap
}
