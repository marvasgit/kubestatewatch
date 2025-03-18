package config

// Handler contains handler configuration
type Handler struct {
	Slack        Slack
	SlackWebhook SlackWebhook
	Hipchat      Hipchat
	Mattermost   Mattermost
	Flock        Flock
	Webhook      Webhook
	CloudEvent   CloudEvent
	MSTeams      MSTeams
	SMTP         SMTP
	Lark         Lark
	Discord      Discord
	Telegram     Telegram
}

// Resource contains resource configuration
type Resource struct {
	Deployment            ResourceConfig
	ReplicationController ResourceConfig
	ReplicaSet            ResourceConfig
	DaemonSet             ResourceConfig
	StatefulSet           ResourceConfig
	Services              ResourceConfig
	Pod                   ResourceConfig
	Job                   ResourceConfig
	Node                  ResourceConfig
	ClusterRole           ResourceConfig
	ClusterRoleBinding    ResourceConfig
	ServiceAccount        ResourceConfig
	PersistentVolume      ResourceConfig
	Namespace             ResourceConfig
	Secret                ResourceConfig
	ConfigMap             ResourceConfig
	Ingress               ResourceConfig
	HPA                   ResourceConfig
	Event                 ResourceConfig
	CoreEvent             ResourceConfig
}

// Config struct contains statemonitor configuration
type Config struct {
	// Handlers know how to send notifications to specific services.
	Handler Handler

	// Resources to watch.
	Resource Resource

	// Configurations for namespaces ot watch or ignore
	NamespacesConfig NamespacesConfig
	// Message properties .
	Message Message
	// Diff properties .
	Diff Diff
}

type NamespacesConfig struct {
	// For watching specific namespaces, leave it empty for watching all.
	// this config is ignored when watching namespaces as resource
	Include []string
	// For ignoring specific namespaces
	Exclude []string
}

type ResourceConfig struct {
	Enabled bool
	// process events based on its type
	//create, update, delete
	//if empty, all events will be processed
	IncludeEvenTypes []string
	IgnorePath       []string
}
type Diff struct {
	//IgnorePath for all resources
	IgnorePath []string
}

// Message contains message configuration.
type Message struct {
	// Message title.
	Title string
}

type Discord struct {
	Enabled    bool
	WebhookURL string
}

// Slack contains slack configuration
type Slack struct {
	// Enable slack notifications.
	Enabled bool
	// Slack "legacy" API token.
	Token string
	// Slack channel.
	Channel string
	// Title of the message.
	//Title string  // moved to Message
}

// SlackWebhook contains slack configuration
type SlackWebhook struct {
	Enabled bool
	// Slack channel.
	Channel string
	// Slack Username.
	Username string
	// Slack Emoji.
	Emoji string
	// Slack Webhook Url.
	Slackwebhookurl string
}

// Hipchat contains hipchat configuration
type Hipchat struct {
	Enabled bool
	// Hipchat token.
	Token string
	// Room name.
	Room string
	// URL of the hipchat server.
	Url string
}

// Mattermost contains mattermost configuration
type Mattermost struct {
	Enabled  bool
	Channel  string
	Url      string
	Username string
}

// Flock contains flock configuration
type Flock struct {
	Enabled bool
	// URL of the flock API.
	Url string
}

// Webhook contains webhook configuration
type Webhook struct {
	Enabled bool
	// Webhook URL.
	Url     string
	Cert    string
	TlsSkip bool
}

// Lark contains lark configuration
type Lark struct {
	Enabled bool
	// Webhook URL.
	WebhookURL string
}

// CloudEvent contains CloudEvent configuration
type CloudEvent struct {
	Enabled bool
	Url     string
}

// MSTeams contains MSTeams configuration
type MSTeams struct {
	Enabled bool
	// MSTeams API Webhook URL.
	WebhookURL string
}

// SMTP contains SMTP configuration.
type SMTP struct {
	Enabled bool
	// Destination e-mail address.
	To string
	// Sender e-mail address .
	From string
	// Smarthost, aka "SMTP server"; address of server used to send email.
	Smarthost string
	// Subject of the outgoing emails.
	Subject string
	// Extra e-mail headers to be added to all outgoing messages.
	Headers map[string]string
	// Authentication parameters.
	Auth SMTPAuth
	// If "true" forces secure SMTP protocol (AKA StartTLS).
	RequireTLS bool
	// SMTP hello field (optional)
	Hello string
}

type SMTPAuth struct {
	Enabled bool
	// Username for PLAN and LOGIN auth mechanisms.
	Username string
	// Password for PLAIN and LOGIN auth mechanisms.
	Password string
	// Identity for PLAIN auth mechanism
	Identity string
	// Secret for CRAM-MD5 auth mechanism
	Secret string
}

// Telegram contains telegram bot configuration
type Telegram struct {
	Enabled bool

	Token           string
	ChatID          int64
	MessageThreadID int64
}
