<div align="center">

**This porject started as fork of kubewatch project maintained by [Robusta.dev](https://home.robusta.dev/) [originally by Bitnami](https://github.com/bitnami-labs/kubewatch/),**

There is way too many changes happening in a Kubernetes cluster, and it is not always easy to keep track of them. diffwatcher is a Kubernetes watcher that publishes notifications to available collaboration hubs/notification channels.When there is an update on any of the watched components there is a diff as notification. Run it in your k8s cluster, and you will get event notifications through webhooks. Because of the diff nature of the notifications, you can easily see what has changed. You can also use it to watch for new resources and get notified when they are created. Because of k8s nature and its regular updates, there is a possiblity to ignore some of the changes, like metadata changes, status changes, etc. This is configurable in the config file.

<img src="./docs/diffwatcher-logo.jpeg">

[![Build Status](https://travis-ci.org/marvasgit/kubernetes-diffwatcher.svg?branch=master)](https://travis-ci.org/marvasgit/kubernetes-diffwatcher) 
[![Go Report Card](https://goreportcard.com/badge/github.com/marvasgit/kubernetes-diffwatcher)](https://goreportcard.com/report/github.com/marvasgit/kubernetes-diffwatcher) 
[![codecov](https://codecov.io/gh/marvasgit/kubernetes-diffwatcher/branch/master/graph/badge.svg)](https://codecov.io/gh/marvasgit/kubernetes-diffwatcher)
[![Docker Pulls](https://img.shields.io/docker/pulls/marvasgit/kubernetes-diffwatcher.svg)](https://hub.docker.com/r/marvasgit/kubernetes-diffwatcher/) 
![GitHub release](https://img.shields.io/github/release/marvasgit/kubernetes-diffwatcher.svg)
</div>

There are basically two kind of notifications:
**notifications for UPDATED items** The whole idea behind is to track the **usefull** differences made on the items we watch, ignoring things like metadata changes, status changes, etc. Not only a simple msg that something was changed.
**notifications for ADDED/DELETED items** this is the original idea behind kubewatch, to track the added/deleted items and notify about them.


# Latest image

```
docmarr/kubernetes-diffwatcher:1.0.0
```

# Install

```console
$ helm repo add diffwatcher https://marvasgit.github.io/kubernetes-diffwatcher/
$ helm install my-release diffwatcher
```

## Introduction

This chart bootstraps a diffwatcher deployment on a [Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```console
$ helm repo add diffwatcher https://marvasgit.github.io/kubernetes-diffwatcher/
$ helm install my-release diffwatcher
```

The command deploys diffwatcher on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

### Global parameters

| Name                      | Description                                     | Value |
| ------------------------- | ----------------------------------------------- | ----- |
| `global.imageRegistry`    | Global Docker image registry                    | `""`  |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]`  |


### Common parameters

| Name                     | Description                                                                             | Value          |
| ------------------------ | --------------------------------------------------------------------------------------- | -------------- |
| `kubeVersion`            | Force target Kubernetes version (using Helm capabilities if not set)                    | `""`           |
| `nameOverride`           | String to partially override common.names.fullname template                             | `""`           |
| `fullnameOverride`       | String to fully override common.names.fullname template                                 | `""`           |
| `commonLabels`           | Labels to add to all deployed objects                                                   | `{}`           |
| `commonAnnotations`      | Annotations to add to all deployed objects                                              | `{}`           |
| `diagnosticMode.enabled` | Enable diagnostic mode (all probes will be disabled and the command will be overridden) | `false`        |
| `diagnosticMode.command` | Command to override all containers in the the deployment(s)/statefulset(s)              | `["sleep"]`    |
| `diagnosticMode.args`    | Args to override all containers in the the deployment(s)/statefulset(s)                 | `["infinity"]` |
| `extraDeploy`            | Array of extra objects to deploy with the release                                       | `[]`           |


### diffwatcher parameters

| Name                                     | Description                                                                      | Value                  |
| ---------------------------------------- | -------------------------------------------------------------------------------- | ---------------------- |
| `image.registry`                         | diffwatcher image registry                                                         | `docker.io`            |
| `image.repository`                       | diffwatcher image repository                                                       | `bitnami/diffwatcher`    |
| `image.tag`                              | diffwatcher image tag (immutable tags are recommended)                             | `0.1.0-debian-10-r513` |
| `image.pullPolicy`                       | diffwatcher image pull policy                                                      | `IfNotPresent`         |
| `image.pullSecrets`                      | Specify docker-registry secret names as an array                                 | `[]`                   |
| `hostAliases`                            | Add deployment host aliases                                                      | `[]`                   |
| `slack.enabled`                          | Enable Slack notifications                                                       | `true`                 |
| `slack.channel`                          | Slack channel to notify                                                          | `XXXX`                 |
| `slack.token`                            | Slack API token                                                                  | `XXXX`                 |
| `hipchat.enabled`                        | Enable HipChat notifications                                                     | `false`                |
| `hipchat.room`                           | HipChat room to notify                                                           | `""`                   |
| `hipchat.token`                          | HipChat token                                                                    | `""`                   |
| `hipchat.url`                            | HipChat URL                                                                      | `""`                   |
| `mattermost.enabled`                     | Enable Mattermost notifications                                                  | `false`                |
| `mattermost.channel`                     | Mattermost channel to notify                                                     | `""`                   |
| `mattermost.url`                         | Mattermost URL                                                                   | `""`                   |
| `mattermost.username`                    | Mattermost user to notify                                                        | `""`                   |
| `flock.enabled`                          | Enable Flock notifications                                                       | `false`                |
| `flock.url`                              | Flock URL                                                                        | `""`                   |
| `msteams.enabled`                        | Enable Microsoft Teams notifications                                             | `false`                |
| `msteams.webhookurl`                     | Microsoft Teams webhook URL                                                      | `""`                   |
| `webhook.enabled`                        | Enable Webhook notifications                                                     | `false`                |
| `webhook.url`                            | Webhook URL                                                                      | `""`                   |
| `smtp.enabled`                           | Enable SMTP (email) notifications                                                | `false`                |
| `smtp.to`                                | Destination email address (required)                                             | `""`                   |
| `smtp.from`                              | Source email address (required)                                                  | `""`                   |
| `smtp.hello`                             | SMTP hello field (optional)                                                      | `""`                   |
| `smtp.smarthost`                         | SMTP server address (name:port) (required)                                       | `""`                   |
| `smtp.subject`                           | Source email subject                                                             | `""`                   |
| `smtp.auth.username`                     | Username for LOGIN and PLAIN auth mech                                           | `""`                   |
| `smtp.auth.password`                     | Password for LOGIN and PLAIN auth mech                                           | `""`                   |
| `smtp.auth.secret`                       | Secret for CRAM-MD5 auth mech                                                    | `""`                   |
| `smtp.auth.identity`                     | Identity for PLAIN auth mech                                                     | `""`                   |
| `smtp.requireTLS`                        | Force STARTTLS. Set to `true` or `false`                                         | `""`                   |
| `namespaceToWatch`                       | Namespace to watch, leave it empty for watching all                              | `""`                   |
| `resourcesToWatch.deployment`            | Watch changes to Deployments                                                     | `true`                 |
| `resourcesToWatch.replicationcontroller` | Watch changes to ReplicationControllers                                          | `false`                |
| `resourcesToWatch.replicaset`            | Watch changes to ReplicaSets                                                     | `false`                |
| `resourcesToWatch.daemonset`             | Watch changes to DaemonSets                                                      | `false`                |
| `resourcesToWatch.services`              | Watch changes to Services                                                        | `false`                |
| `resourcesToWatch.pod`                   | Watch changes to Pods                                                            | `true`                 |
| `resourcesToWatch.job`                   | Watch changes to Jobs                                                            | `false`                |
| `resourcesToWatch.persistentvolume`      | Watch changes to PersistentVolumes                                               | `false`                |
| `command`                                | Override default container command (useful when using custom images)             | `[]`                   |
| `args`                                   | Override default container args (useful when using custom images)                | `[]`                   |
| `lifecycleHooks`                         | for the diffwatcher container(s) to automate configuration before or after startup | `{}`                   |
| `extraEnvVars`                           | Extra environment variables to be set on diffwatcher container                     | `[]`                   |
| `extraEnvVarsCM`                         | Name of existing ConfigMap containing extra env vars                             | `""`                   |
| `extraEnvVarsSecret`                     | Name of existing Secret containing extra env vars                                | `""`                   |


### diffwatcher deployment parameters

| Name                                    | Description                                                                               | Value           |
| --------------------------------------- | ----------------------------------------------------------------------------------------- | --------------- |
| `replicaCount`                          | Number of diffwatcher replicas to deploy                                                    | `1`             |
| `podSecurityContext.enabled`            | Enable diffwatcher containers' SecurityContext                                              | `false`         |
| `podSecurityContext.fsGroup`            | Set diffwatcher containers' SecurityContext fsGroup                                         | `""`            |
| `containerSecurityContext.enabled`      | Enable diffwatcher pods' Security Context                                                   | `false`         |
| `containerSecurityContext.runAsUser`    | Set diffwatcher pods' SecurityContext runAsUser                                             | `""`            |
| `containerSecurityContext.runAsNonRoot` | Set diffwatcher pods' SecurityContext runAsNonRoot                                          | `""`            |
| `resources.limits`                      | The resources limits for the diffwatcher container                                          | `{}`            |
| `resources.requests`                    | The requested resources for the diffwatcher container                                       | `{}`            |
| `startupProbe.enabled`                  | Enable startupProbe                                                                       | `false`         |
| `startupProbe.initialDelaySeconds`      | Initial delay seconds for startupProbe                                                    | `10`            |
| `startupProbe.periodSeconds`            | Period seconds for startupProbe                                                           | `10`            |
| `startupProbe.timeoutSeconds`           | Timeout seconds for startupProbe                                                          | `1`             |
| `startupProbe.failureThreshold`         | Failure threshold for startupProbe                                                        | `3`             |
| `startupProbe.successThreshold`         | Success threshold for startupProbe                                                        | `1`             |
| `livenessProbe.enabled`                 | Enable livenessProbe                                                                      | `false`         |
| `livenessProbe.initialDelaySeconds`     | Initial delay seconds for livenessProbe                                                   | `10`            |
| `livenessProbe.periodSeconds`           | Period seconds for livenessProbe                                                          | `10`            |
| `livenessProbe.timeoutSeconds`          | Timeout seconds for livenessProbe                                                         | `1`             |
| `livenessProbe.failureThreshold`        | Failure threshold for livenessProbe                                                       | `3`             |
| `livenessProbe.successThreshold`        | Success threshold for livenessProbe                                                       | `1`             |
| `readinessProbe.enabled`                | Enable readinessProbe                                                                     | `false`         |
| `readinessProbe.initialDelaySeconds`    | Initial delay seconds for readinessProbe                                                  | `10`            |
| `readinessProbe.periodSeconds`          | Period seconds for readinessProbe                                                         | `10`            |
| `readinessProbe.timeoutSeconds`         | Timeout seconds for readinessProbe                                                        | `1`             |
| `readinessProbe.failureThreshold`       | Failure threshold for readinessProbe                                                      | `3`             |
| `readinessProbe.successThreshold`       | Success threshold for readinessProbe                                                      | `1`             |
| `customStartupProbe`                    | Override default startup probe                                                            | `{}`            |
| `customLivenessProbe`                   | Override default liveness probe                                                           | `{}`            |
| `customReadinessProbe`                  | Override default readiness probe                                                          | `{}`            |
| `podAffinityPreset`                     | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`       | `""`            |
| `podAntiAffinityPreset`                 | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`  | `soft`          |
| `nodeAffinityPreset.type`               | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard` | `""`            |
| `nodeAffinityPreset.key`                | Node label key to match. Ignored if `affinity` is set.                                    | `""`            |
| `nodeAffinityPreset.values`             | Node label values to match. Ignored if `affinity` is set.                                 | `[]`            |
| `affinity`                              | Affinity for pod assignment                                                               | `{}`            |
| `nodeSelector`                          | Node labels for pod assignment                                                            | `{}`            |
| `tolerations`                           | Tolerations for pod assignment                                                            | `[]`            |
| `priorityClassName`                     | Controller priorityClassName                                                              | `""`            |
| `schedulerName`                         | Name of the k8s scheduler (other than default)                                            | `""`            |
| `topologySpreadConstraints`             | Topology Spread Constraints for pod assignment                                            | `[]`            |
| `podLabels`                             | Extra labels for diffwatcher pods                                                           | `{}`            |
| `podAnnotations`                        | Annotations for diffwatcher pods                                                            | `{}`            |
| `extraVolumes`                          | Optionally specify extra list of additional volumes for diffwatcher pods                    | `[]`            |
| `extraVolumeMounts`                     | Optionally specify extra list of additional volumeMounts for diffwatcher container(s)       | `[]`            |
| `updateStrategy.type`                   | Deployment strategy type.                                                                 | `RollingUpdate` |
| `initContainers`                        | Add additional init containers to the diffwatcher pods                                      | `[]`            |
| `sidecars`                              | Add additional sidecar containers to the diffwatcher pods                                   | `[]`            |


### RBAC parameters

| Name                                          | Description                                                                                                         | Value   |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------- |
| `rbac.create`                                 | Whether to create & use RBAC resources or not                                                                       | `false` |
| `serviceAccount.create`                       | Specifies whether a ServiceAccount should be created                                                                | `true`  |
| `serviceAccount.name`                         | Name of the service account to use. If not set and create is true, a name is generated using the fullname template. | `""`    |
| `serviceAccount.automountServiceAccountToken` | Automount service account token for the server service account                                                      | `true`  |
| `serviceAccount.annotations`                  | Annotations for service account. Evaluated as a template. Only used if `create` is `true`.                          | `{}`    |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
$ helm install my-release bitnami/diffwatcher \
  --set=slack.channel="#bots",slack.token="XXXX-XXXX-XXXX"
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install my-release -f values.yaml bitnami/diffwatcher
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Configuration and installation details

### [Rolling VS Immutable tags](https://docs.bitnami.com/containers/how-to/understand-rolling-tags-containers/)

It is strongly recommended to use immutable tags in a production environment. This ensures your deployment does not change automatically if the same tag is updated with a different image.

Bitnami will release a new chart updating its containers if a new version of the main container, significant changes, or critical vulnerabilities exist.

### Create a Slack bot

Open [https://my.slack.com/services/new/bot](https://my.slack.com/services/new/bot) to create a new Slack bot.
The API token can be found on the edit page (it starts with `xoxb-`).

Invite the Bot to your channel by typing `/join @name_of_your_bot` in the Slack message area.

### Adding extra environment variables

In case you want to add extra environment variables (useful for advanced operations like custom init scripts), you can use the `extraEnvVars` property.

```yaml
extraEnvVars:
  - name: LOG_LEVEL
    value: debug
```

Alternatively, you can use a ConfigMap or a Secret with the environment variables. To do so, use the `extraEnvVarsCM` or the `extraEnvVarsSecret` values.

### Sidecars and Init Containers

If you have a need for additional containers to run within the same pod as the diffwatcher app (e.g. an additional metrics or logging exporter), you can do so via the `sidecars` config parameter. Simply define your container according to the Kubernetes container spec.

```yaml
sidecars:
  - name: your-image-name
    image: your-image
    imagePullPolicy: Always
    ports:
      - name: portname
       containerPort: 1234
```

Similarly, you can add extra init containers using the `initContainers` parameter.

```yaml
initContainers:
  - name: your-image-name
    image: your-image
    imagePullPolicy: Always
    ports:
      - name: portname
        containerPort: 1234
```

### Deploying extra resources

There are cases where you may want to deploy extra objects, such a ConfigMap containing your app's configuration or some extra deployment with a micro service used by your app. For covering this case, the chart allows adding the full specification of other objects using the `extraDeploy` parameter.

### Setting Pod's affinity

This chart allows you to set your custom affinity using the `affinity` parameter. Find more information about Pod's affinity in the [kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity).

As an alternative, you can use of the preset configurations for pod affinity, pod anti-affinity, and node affinity available at the [bitnami/common](https://github.com/bitnami/charts/tree/master/bitnami/common#affinities) chart. To do so, set the `podAffinityPreset`, `podAntiAffinityPreset`, or `nodeAffinityPreset` parameters.


### Local Installation
#### Using go package installer:

```console
# Download and install diffwatcher
$ go get -u github.com/marvasgit/kubernetes-diffwatcher

# Configure the notification channel
$ diffwatcher config add slack --channel <slack_channel> --token <slack_token>

# Add resources to be watched
$ diffwatcher resource add --po --svc
INFO[0000] resource svc configured
INFO[0000] resource po configured

# start diffwatcher server
$ diffwatcher
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-service
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-pod
INFO[0000] Processing add to service: default/kubernetes  pkg=diffwatcher-service
INFO[0000] Processing add to service: kube-system/tiller-deploy  pkg=diffwatcher-service
INFO[0000] Processing add to pod: kube-system/tiller-deploy-69ffbf64bc-h8zxm  pkg=diffwatcher-pod
INFO[0000] diffwatcher controller synced and ready         pkg=diffwatcher-service
INFO[0000] diffwatcher controller synced and ready         pkg=diffwatcher-pod

```
#### Using Docker:

To Run diffwatcher Container interactively, place the config file in `$HOME/.diffwatcher.yaml` location and use the following command.

```
docker run --rm -it --network host -v $HOME/.diffwatcher.yaml:/root/.diffwatcher.yaml -v $HOME/.kube/config:/opt/bitnami/diffwatcher/.kube/config --name <container-name> us-central1-docker.pkg.dev/genuine-flight-317411/devel/diffwatcher
```

Example:

```
$ docker run --rm -it --network host -v $HOME/.diffwatcher.yaml:/root/.diffwatcher.yaml -v $HOME/.kube/config:/opt/bitnami/diffwatcher/.kube/config --name diffwatcher-app us-central1-docker.pkg.dev/genuine-flight-317411/devel/diffwatcher

==> Writing config file...
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-service
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-pod
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-deployment
INFO[0000] Starting diffwatcher controller                 pkg=diffwatcher-namespace
INFO[0000] Processing add to namespace: kube-node-lease  pkg=diffwatcher-namespace
INFO[0000] Processing add to namespace: kube-public      pkg=diffwatcher-namespace
INFO[0000] Processing add to namespace: kube-system      pkg=diffwatcher-namespace
INFO[0000] Processing add to namespace: default          pkg=diffwatcher-namespace
....
```

To Demonise diffwatcher container use

```
$ docker run --rm -d --network host -v $HOME/.diffwatcher.yaml:/root/.diffwatcher.yaml -v $HOME/.kube/config:/opt/bitnami/diffwatcher/.kube/config --name diffwatcher-app us-central1-docker.pkg.dev/genuine-flight-317411/devel/diffwatcher
```

# Configure

diffwatcher supports `config` command for configuration. Config file will be saved at `$HOME/.diffwatcher.yaml`

```
$ diffwatcher config -h

config command allows admin setup his own configuration for running diffwatcher

Usage:
  diffwatcher config [flags]
  diffwatcher config [command]

Available Commands:
  add         add webhook config to .diffwatcher.yaml
  test        test handler config present in .diffwatcher.yaml
  view        view .diffwatcher.yaml

Flags:
  -h, --help   help for config

Use "diffwatcher config [command] --help" for more information about a command.
```
### Example:

### slack:

- Create a [slack Bot](https://my.slack.com/services/new/bot)

- Edit the Bot to customize its name, icon and retrieve the API token (it starts with `xoxb-`).

- Invite the Bot into your channel by typing: `/invite @name_of_your_bot` in the Slack message area.

- Add Api token to diffwatcher config using the following steps

  ```console
  $ diffwatcher config add slack --channel <slack_channel> --token <slack_token>
  ```
  You have an altenative choice to set your SLACK token, channel via environment variables:

  ```console
  $ export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
  $ export KW_SLACK_CHANNEL='#channel_name'
  ```

### slackwebhookurl:

- Create a [slack app](https://api.slack.com/apps/new)

- Enable [Incoming Webhooks](https://api.slack.com/messaging/webhooks#enable_webhooks). (On "Settings" page.)

- Create an incoming webhook URL (Add New Webhook to Workspace on "Settings" page.)

- Pick a channel that the app will post to, and then click to Authorize your app. You will get back your webhook URL.  
  The Slack Webhook URL will look like: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

- Add slack webhook url to diffwatcher config using the following steps

  ```console
  $ diffwatcher config add slackwebhookurl --username <slack_username> --emoji <slack_emoji> --channel <slack_channel> --slackwebhookurl <slack_webhook_url>
  ```
  Or, you have an altenative choice to set your SLACK channel, username, emoji and webhook URL via environment variables:

  ```console
  $ export KW_SLACK_CHANNEL=slack_channel
  $ export KW_SLACK_USERNAME=slack_username
  $ export KW_SLACK_EMOJI=slack_emoji
  $ export KW_SLACK_WEBHOOK_URL=slack_webhook_url
  ```
  
 - Example apply done in a bash script:  
  
 ```console
 $ cat diffwatcher-configmap-slackwebhook.yaml | sed "s|<slackchannel>|"\"$SlackChannel"\"|g;s|<slackusername>|"\"$SlackUsesrName"\"|g;s|<slackemoji>|"\"$SlackEmoji"\"|g;s|<SlackWebhookUrl>|"\"$WebhookUrl"\"|g" | kubectl create -f -
 ```
 
 - An example diffwatcher-configmap-slackwebhook.yaml YAML File:  
  
 ```yaml
 apiVersion: v1
kind: ConfigMap
metadata:
  name: diffwatcher
data:
  .diffwatcher.yaml: |
    namespace: ""
    handler:
      slackwebhook:
        enabled: true
        channel: <slackchannel>
        username: <slackusername>
        emoji: <slackemoji>
        slackwebhookurl: <SlackWebhookUrl>
    resource:
      clusterrole: false
      configmap: false
      daemonset: false
      deployment: true
      ingress: false
      job: false
      namespace: false
      node: false
      persistentvolume: false
      pod: true
      replicaset: false
      replicationcontroller: false
      secret: false
      serviceaccount: false
      services: true
      event: true
      coreevent: false
    ```
```
### flock:

- Create a [flock bot](https://docs.flock.com/display/flockos/Bots).

- Add flock webhook url to config using the following command.
  ```console
  $ diffwatcher config add flock --url <flock_webhook_url>
  ```
  You have an altenative choice to set your FLOCK URL

  ```console
  $ export KW_FLOCK_URL='https://api.flock.com/hooks/sendMessage/XXXXXXXX'
  ```

## Testing Config

To test the handler config by send test messages use the following command.
```
$ diffwatcher config test -h

Tests handler configs present in .diffwatcher.yaml by sending test messages

Usage:
  diffwatcher config test [flags]

Flags:
  -h, --help   help for test
```

#### Example:

```
$ diffwatcher config test

Testing Handler configs from .diffwatcher.yaml
2019/06/03 12:29:23 Message successfully sent to channel ABCD at 1559545162.000100
```

## Viewing config
To view the entire config file `$HOME/.diffwatcher.yaml` use the following command.
```
$ diffwatcher config view
Contents of .diffwatcher.yaml

handler:
  slack:
    token: xoxb-xxxxx-yyyy-zzz
    channel: kube-watch
  hipchat:
    token: ""
    room: ""
    url: ""
  mattermost:
    channel: ""
    url: ""
    username: ""
  flock:
    url: ""
  webhook:
    url: ""
  cloudevent:
    url: ""
resource:
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: false
  pod: true
  job: false
  node: false
  clusterrole: false
  clusterrolebinding: false
  serviceaccount: false
  persistentvolume: false
  namespace: false
  secret: false
  configmap: false
  ingress: false
  event: true
  coreevent: false
namespace: ""
message:
  title: "XXXX"
diff:
  ignore:
  - "metadata.creationTimestamp"
  - "metadata.resourceVersion"
  - "metadata.selfLink"
  - "metadata.uid"
  - "status"
  - "metadata.managedFields"
  - "metadata.generation"
  - "metadata.annotations.kubectl.kubernetes.io/last-applied-configuration"
  - "metadata.annotations.deployment.kubernetes.io/revision"
  - "metadata.annotations.kubernetes.io/change-cause"
  - "metadata.annotations.kubernetes.io/psp"
  - "metadata.annotations.kubernetes.io/psp-spec"
  - "metadata.annotations.kubernetes.io/psp-status"

```


## Resources

To manage the resources being watched, use the following command, changes will be saved to `$HOME/.diffwatcher.yaml`.

```
$ diffwatcher resource -h

manage resources to be watched

Usage:
  diffwatcher resource [flags]
  diffwatcher resource [command]

Available Commands:
  add         adds specific resources to be watched
  remove      remove specific resources being watched

Flags:
      
      --clusterrolebinding      watch for cluster role bindings
      --clusterrole             watch for cluster roles
      --cm                      watch for plain configmaps
      --deploy                  watch for deployments
      --ds                      watch for daemonsets
  -h, --help                    help for resource
      --ing                     watch for ingresses
      --job                     watch for jobs
      --node                    watch for Nodes
      --ns                      watch for namespaces
      --po                      watch for pods
      --pv                      watch for persistent volumes
      --rc                      watch for replication controllers
      --rs                      watch for replicasets
      --sa                      watch for service accounts
      --secret                  watch for plain secrets
      --svc                     watch for services
      --coreevent               watch for events from the kubernetes core api. (Old events api, replaced in kubernetes 1.19)

Use "diffwatcher resource [command] --help" for more information about a command.

```

### Add/Remove resource:
```
$ diffwatcher resource add -h

adds specific resources to be watched

Usage:
  diffwatcher resource add [flags]

Flags:
  -h, --help   help for add

Global Flags:
      --clusterrole             watch for cluster roles
      --clusterrolebinding      watch for cluster role bindings
      --cm                      watch for plain configmaps
      --deploy                  watch for deployments
      --ds                      watch for daemonsets
      --ing                     watch for ingresses
      --job                     watch for jobs
      --node                    watch for Nodes
      --ns                      watch for namespaces
      --po                      watch for pods
      --pv                      watch for persistent volumes
      --rc                      watch for replication controllers
      --rs                      watch for replicasets
      --sa                      watch for service accounts
      --secret                  watch for plain secrets
      --svc                     watch for services
      --coreevent               watch for events from the kubernetes core api. (Old events api, replaced in kubernetes 1.19)

```

### Example:

```console
# rc, po and svc will be watched
$ diffwatcher resource add --rc --po --svc

# rc, po and svc will be stopped from being watched
$ diffwatcher resource remove --rc --po --svc
```

### Changing log level

In case you want to change the default log level, add an environment variable named `LOG_LEVEL` with value from `trace/debug/info/warning/error` 

```yaml
env:
- name: LOG_LEVEL
  value: debug
```

### Changing log format

In case you want to change the log format to `json`, add an environment variable named `LOG_FORMATTER` with value `json`

```yaml
env:
- name: LOG_FORMATTER
  value: json
```

# Build

### Using go

Clone the repository anywhere:
```console
$ git clone https://github.com/marvasgit/kubernetes-diffwatcher.git
$ cd diffwatcher
$ go build
```
or

You can also use the Makefile directly:

```console
$ make build
```

#### Prerequisites

- You need to have [Go](http://golang.org) (v1.5 or later)  installed. Make sure to set `$GOPATH`


### Using Docker

```console
$ make docker-image
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
diffwatcher           latest              919896d3cd90        3 minutes ago       27.9MB
```
#### Prerequisites

- you need to have [docker](https://docs.docker.com/) installed.

# Contribution

Refer to the [contribution guidelines](docs/CONTRIBUTION.md) to get started.
