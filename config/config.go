package config

// Handler contains handler configuration
type Handler struct {
	Slack        Slack        `yaml:"slack"`
	SlackWebhook SlackWebhook `yaml:"slackwebhook"`
	Hipchat      Hipchat      `yaml:"hipchat"`
	Mattermost   Mattermost   `yaml:"mattermost"`
	Flock        Flock        `yaml:"flock"`
	Webhook      Webhook      `yaml:"webhook"`
	CloudEvent   CloudEvent   `yaml:"cloudevent"`
	MSTeams      MSTeams      `yaml:"msteams"`
	SMTP         SMTP         `yaml:"smtp"`
	Lark         Lark         `yaml:"lark"`
}

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `yaml:"deployment"`
	ReplicationController bool `yaml:"rc"`
	ReplicaSet            bool `yaml:"rs"`
	DaemonSet             bool `yaml:"ds"`
	StatefulSet           bool `yaml:"statefulset"`
	Services              bool `yaml:"svc"`
	Pod                   bool `yaml:"po"`
	Job                   bool `yaml:"job"`
	Node                  bool `yaml:"node"`
	ClusterRole           bool `yaml:"clusterrole"`
	ClusterRoleBinding    bool `yaml:"clusterrolebinding"`
	ServiceAccount        bool `yaml:"sa"`
	PersistentVolume      bool `yaml:"pv"`
	Namespace             bool `yaml:"ns"`
	Secret                bool `yaml:"secret"`
	ConfigMap             bool `yaml:"configmap"`
	Ingress               bool `yaml:"ing"`
	HPA                   bool `yaml:"hpa"`
	Event                 bool `yaml:"event"`
	CoreEvent             bool `yaml:"coreevent"`
}

// Config struct contains diffwatcher configuration
type Config struct {
	// Handlers know how to send notifications to specific services.
	Handler Handler `yaml:"handler"`

	//Reason   []string `yaml:"reason"`

	// Resources to watch.
	Resource Resource `yaml:"resource"`

	// Configurations for namespaces ot watch or ignore
	NamespacesConfig NamespacesConfig `yaml:"namespacesconfig"`
	// Message properties .
	Message Message `yaml:"message"`
	// Diff properties .
	Diff Diff `yaml:"diff"`
}

type NamespacesConfig struct {
	// For watching specific namespaces, leave it empty for watching all.
	// this config is ignored when watching namespaces as resource
	Include []string `yaml:"include"`
	// For ignoring specific namespaces
	Exclude []string `yaml:"exclude"`
}

type Diff struct {
	IgnorePath []string `yaml:"ignore"`
}

// Message contains message configuration.
type Message struct {
	// Message title.
	Title string `yaml:"title"`
}

// Slack contains slack configuration
type Slack struct {
	// Enable slack notifications.
	Enabled bool `yaml:"enabled"`
	// Slack "legacy" API token.
	Token string `yaml:"token"`
	// Slack channel.
	Channel string `yaml:"channel"`
	// Title of the message.
	//Title string `yaml:"title"` // moved to Message
}

// SlackWebhook contains slack configuration
type SlackWebhook struct {
	Enabled bool `yaml:"enabled"`
	// Slack channel.
	Channel string `yaml:"channel"`
	// Slack Username.
	Username string `yaml:"username"`
	// Slack Emoji.
	Emoji string `yaml:"emoji"`
	// Slack Webhook Url.
	Slackwebhookurl string `yaml:"slackwebhookurl"`
}

// Hipchat contains hipchat configuration
type Hipchat struct {
	Enabled bool `yaml:"enabled"`
	// Hipchat token.
	Token string `yaml:"token"`
	// Room name.
	Room string `yaml:"room"`
	// URL of the hipchat server.
	Url string `yaml:"url"`
}

// Mattermost contains mattermost configuration
type Mattermost struct {
	Enabled  bool   `yaml:"enabled"`
	Channel  string `yaml:"room"`
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
}

// Flock contains flock configuration
type Flock struct {
	Enabled bool `yaml:"enabled"`
	// URL of the flock API.
	Url string `yaml:"url"`
}

// Webhook contains webhook configuration
type Webhook struct {
	Enabled bool `yaml:"enabled"`
	// Webhook URL.
	Url     string `yaml:"url"`
	Cert    string `yaml:"cert"`
	TlsSkip bool   `yaml:"tlsskip"`
}

// Lark contains lark configuration
type Lark struct {
	Enabled bool `yaml:"enabled"`
	// Webhook URL.
	WebhookURL string `yaml:"webhookurl"`
}

// CloudEvent contains CloudEvent configuration
type CloudEvent struct {
	Enabled bool   `yaml:"enabled"`
	Url     string `yaml:"url"`
}

// MSTeams contains MSTeams configuration
type MSTeams struct {
	Enabled bool `yaml:"enabled"`
	// MSTeams API Webhook URL.
	WebhookURL string `yaml:"webhookurl"`
}

// SMTP contains SMTP configuration.
type SMTP struct {
	Enabled bool `yaml:"enabled"`
	// Destination e-mail address.
	To string `yaml:"to" yaml:"to,omitempty"`
	// Sender e-mail address .
	From string `yaml:"from" yaml:"from,omitempty"`
	// Smarthost, aka "SMTP server"; address of server used to send email.
	Smarthost string `yaml:"smarthost" yaml:"smarthost,omitempty"`
	// Subject of the outgoing emails.
	Subject string `yaml:"subject" yaml:"subject,omitempty"`
	// Extra e-mail headers to be added to all outgoing messages.
	Headers map[string]string `yaml:"headers" yaml:"headers,omitempty"`
	// Authentication parameters.
	Auth SMTPAuth `yaml:"auth" yaml:"auth,omitempty"`
	// If "true" forces secure SMTP protocol (AKA StartTLS).
	RequireTLS bool `yaml:"requireTLS" yaml:"requireTLS"`
	// SMTP hello field (optional)
	Hello string `yaml:"hello" yaml:"hello,omitempty"`
}

type SMTPAuth struct {
	Enabled bool `yaml:"enabled"`
	// Username for PLAN and LOGIN auth mechanisms.
	Username string `yaml:"username" yaml:"username,omitempty"`
	// Password for PLAIN and LOGIN auth mechanisms.
	Password string `yaml:"password" yaml:"password,omitempty"`
	// Identity for PLAIN auth mechanism
	Identity string `yaml:"identity" yaml:"identity,omitempty"`
	// Secret for CRAM-MD5 auth mechanism
	Secret string `yaml:"secret" yaml:"secret,omitempty"`
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
