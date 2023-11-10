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
}

// Resource contains resource configuration
type Resource struct {
	Deployment            bool
	ReplicationController bool
	ReplicaSet            bool
	DaemonSet             bool
	StatefulSet           bool
	Services              bool
	Pod                   bool
	Job                   bool
	Node                  bool
	ClusterRole           bool
	ClusterRoleBinding    bool
	ServiceAccount        bool
	PersistentVolume      bool
	Namespace             bool
	Secret                bool
	ConfigMap             bool
	Ingress               bool
	HPA                   bool
	Event                 bool
	CoreEvent             bool
}

// Config struct contains diffwatcher configuration
type Config struct {
	// Handlers know how to send notifications to specific services.
	Handler Handler

	//Reason   []string

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

type Diff struct {
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

// CheckMissingResourceEnvvars will read the environment for equivalent config variables to set
func (c *Config) CheckMissingResourceEnvvars() {
	if !c.Resource.DaemonSet {
		c.Resource.DaemonSet = true
	}
	if !c.Resource.ReplicaSet {
		c.Resource.ReplicaSet = true
	}
	if !c.Resource.Namespace {
		c.Resource.Namespace = true
	}
	if !c.Resource.Deployment {
		c.Resource.Deployment = true
	}
	if !c.Resource.Pod {
		c.Resource.Pod = true
	}
	if !c.Resource.ReplicationController {
		c.Resource.ReplicationController = true
	}
	if !c.Resource.Services {
		c.Resource.Services = true
	}
	if !c.Resource.Job {
		c.Resource.Job = true
	}
	if !c.Resource.PersistentVolume {
		c.Resource.PersistentVolume = true
	}
	if !c.Resource.Secret {
		c.Resource.Secret = true
	}
	if !c.Resource.ConfigMap {
		c.Resource.ConfigMap = true
	}
	if !c.Resource.Ingress {
		c.Resource.Ingress = true
	}
	if !c.Resource.Node {
		c.Resource.Node = true
	}
	if !c.Resource.ServiceAccount {
		c.Resource.ServiceAccount = true
	}
	if !c.Resource.ClusterRole {
		c.Resource.ClusterRole = true
	}
	if !c.Resource.ClusterRoleBinding {
		c.Resource.ClusterRoleBinding = true
	}
}
