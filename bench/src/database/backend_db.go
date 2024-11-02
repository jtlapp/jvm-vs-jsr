package database

var databaseConfig = DatabaseConfig{
	UrlEnvVar:      "BACKEND_DATABASE_URL",
	UsernameEnvVar: "BACKEND_DATABASE_USERNAME",
	PasswordEnvVar: "BACKEND_DATABASE_PASSWORD",
}

type BackendDB struct {
	Database
}

func NewBackendDatabase() *BackendDB {
	return &BackendDB{*NewDatabase(&databaseConfig)}
}
