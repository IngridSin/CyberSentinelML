package config

const (

	// Database Constants
	DBUser     = "postgres"
	DBPassword = "toor@1234567"
	DBName     = "test"
	DBHost     = "192.168.0.5" // Inside SSH tunnel
	DBPort     = "5432"
	DBSchema   = "test_schema"
	DBTable    = "test_network_flow"
	IMDBName   = "im_db"
	IMSchema   = "public"
	IMTable    = "emails"

	// SSH Tunnel Config
	SSHHost = "4.246.216.81"
	SSHPort = "22"
	SSHUser = "Victim"
	SSHKey  = "/Users/grid/Downloads/victim-key.pem"

	// Flask Config
	FlaskUrl   = "http://127.0.0.1:5000"
	PredictAPI = "/predict"

	// Network Capture
	NetworkInterface = "en0"

	// IMAP Config
	IMAPServer = "imap.gmail.com:993"
	IMUsername = "elvismamickey@gmail.com"
	IMPassword = "gmec zazr oack pijz"

	// Application Constants
	AppVersion = "1.0.0"
	LogFile    = "app.log"
)
