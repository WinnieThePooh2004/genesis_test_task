package settings

type AppSettings struct {
	Port             string
	RatesUrl         string
	Email            *EmailSettings
	ConnectionString string
}

type EmailSettings struct {
	Email    string
	Host     string
	Password string
	Port     int
}
