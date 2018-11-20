package configure

// AppConfig
type AppConfig struct {
	StaticFilePath string
}

// AppProperties :
var AppProperties AppConfig

// SetProperties :
func SetProperties(ap AppConfig) {
	AppProperties = ap
}
