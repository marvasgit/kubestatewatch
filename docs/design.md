# statemonitor

statemonitor contains three components: controller, config, handler

![statemonitor Diagram](statemonitor.png?raw=true "statemonitor Overview")

## Config

The config object contains `statemonitor` configuration, like handlers, filters.

A config object is used to creating new client.

## Controller

The controller initializes using the config object by reading the `.statemonitor.yaml` or command line arguments.
If the parameters are not fully mentioned, the config falls back to read a set of standard environment variables.

Controller creates necessary `SharedIndexInformer`s provided by `kubernetes/client-go` for listening and watching
resource changes. Controller updates this subscription information with Kubernetes API Server.

Whenever, the Kubernetes Controller Manager gets events related to the subscribed resources, it pushes the events to
`SharedIndexInformer`. This in-turn puts the events onto a rate-limiting queue for better handling of the events.

Controller picks the events from the queue and hands over the events to the appropriate handler after
necessary filtering.

## Handler

Handler manages how `statemonitor` handles events.

With each event get from k8s and matched filtering from configuration, it is passed to handler. Currently, `statemonitor` has 8 handlers:

 - `Default`: which just print the event in JSON format
 - `Flock`: which send notification to Flock channel based on information from config
 - `Hipchat`: which send notification to Hipchat room based on information from config
 - `Mattermost`: which send notification to Mattermost channel based on information from config
 - `MS Teams`: which send notification to MS Team incoming webhook based on information from config
 - `Slack`: which send notification to Slack channel based on information from config
 - `Smtp`: which sends notifications to email recipients using a SMTP server obtained from config
 - `Lark`: which sends notifications to Lark incoming webhook based on information from config

More handlers will be added in future.

Each handler must implement the [Handler interface](https://github.com/marvasgit/kubestatewatch/blob/master/pkg/handlers/handler.go#L31)
