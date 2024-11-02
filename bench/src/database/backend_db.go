package database

type BackendDB struct {
	Database
}

func NewBackendDatabase() *BackendDB {
	var databaseConfig = DatabaseConfig{
		UrlEnvVar:      "BACKEND_DATABASE_URL",
		UsernameEnvVar: "BACKEND_DATABASE_USERNAME",
		PasswordEnvVar: "BACKEND_DATABASE_PASSWORD",
	}
	return &BackendDB{*NewDatabase(&databaseConfig)}
}
