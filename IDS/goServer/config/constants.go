package config

const (

	// Database Constants
	DBUser     = "postgres"
	DBPassword = "toor@1234567"
	DBName     = "test"
	DBHost     = "127.0.0.1" // Inside SSH tunnel
	DBPort     = "5432"
	DBSchema   = "test_schema"
	DBTable    = "test_network_flow"

	// SSH Tunnel Config
	SSHHost = "20.169.218.42"
	SSHPort = "22"
	SSHUser = "Victim"
	SSHKey  = "/Users/grid/Downloads/victim-key.pem"

	// Network Capture
	NetworkInterface = "en0"

	// Application Constants
	AppVersion = "1.0.0"
	LogFile    = "app.log"
)
