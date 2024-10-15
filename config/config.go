package config

type configKey string

const (
	Port = configKey("PORT")

	DBDSN = configKey("DB_DSN")

	HttpsHost = configKey("HTTPS_HOST")

	IsSwaggerCreated = configKey("IS_SWAGGER_CREATED")
)

// const (
// 	MaxRoutine           = 20
// 	DefaultServerAddress = "localhost:8080"
// 	TokenExpiration      = 24 * time.Hour
// 	JwtKeySecret         = "very-very-secret-key"
// )

// type ServerConfigData struct {
// 	ServerAddr string
// 	DBDSN      string
// }

// func ConfigInit() error {
// 	flag.StringVar(&ServerConfig.ServerAddr, "a", DefaultServerAddress, "address to run HTTP server on")
// 	flag.StringVar(&ServerConfig.DBDSN, "d", "", "base file path to save URLs")

// 	flag.Parse()

// 	if serverAddr := os.Getenv("RUN_ADDRESS"); serverAddr != "" {
// 		ServerConfig.ServerAddr = serverAddr
// 	}
// 	if dbdsn := os.Getenv("DATABASE_URI"); dbdsn != "" {
// 		ServerConfig.DBDSN = dbdsn
// 	}

// 	ServerConfig.DBDSN = "postgres://postgres:7513@localhost:5432/postgres"
// 	ServerConfig.ServerAddr = DefaultServerAddress

// 	return checkConfig()
// }

// func checkConfig() error {
// 	switch {
// 	case ServerConfig.ServerAddr == "":
// 		return apperror.ErrNoServerAddress
// 	case ServerConfig.DBDSN == "":
// 		return apperror.ErrNoDatabaseDestination
// 	default:
// 		return nil
// 	}
// }

// var ServerConfig ServerConfigData
