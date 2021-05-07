package deployment

import "os"

var env map[string]string

func GetEnv(name string) string {
	return os.Getenv(env[name])
}

func InitiateEnv() {
	switch CurrentRunMode {

	case Test:
		{
			env = map[string]string{
				"SENDGRID_API_KEY":       "SENDGRID_API_KEY",
				"JWT_SECRET_KEY":         "JWT_SECRET_KEY",
				"ELASTIC_INDEX_FOR_TEXT": "ELASTIC_INDEX_FOR_TEXT_TEST",
			}
		}

	case Development:
		{
			env = map[string]string{
				"SENDGRID_API_KEY":       "SENDGRID_API_KEY",
				"JWT_SECRET_KEY":         "JWT_SECRET_KEY",
				"ELASTIC_INDEX_FOR_TEXT": "ELASTIC_INDEX_FOR_TEXT",
			}
		}

	case Production:
		{
			env = map[string]string{
				"SENDGRID_API_KEY":       "SENDGRID_API_KEY",
				"JWT_SECRET_KEY":         "JWT_SECRET_KEY",
				"ELASTIC_INDEX_FOR_TEXT": "ELASTIC_INDEX_FOR_TEXT",
			}
		}
	}
}
