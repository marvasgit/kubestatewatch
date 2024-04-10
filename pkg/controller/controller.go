/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/marvasgit/kubernetes-statemonitor/config"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/event"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/handlers"
	"github.com/marvasgit/kubernetes-statemonitor/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"

	apps_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	events_v1 "k8s.io/api/events/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/strings/slices"
)

const maxRetries = 5
const V1 = "v1"
const AUTOSCALING_V1 = "autoscaling/v1"
const APPS_V1 = "apps/v1"
const BATCH_V1 = "batch/v1"
const RBAC_V1 = "rbac.authorization.k8s.io/v1"
const NETWORKING_V1 = "networking.k8s.io/v1"
const EVENTS_V1 = "events.k8s.io/v1"

var serverStartTime time.Time
var confDiff config.Diff
var namespaces []string
var metric *prometheus.CounterVec
var mu sync.Mutex

// Event indicate the informerEvent
type EventWrapper struct {
	Event          Event
	ResourceConfig *config.ResourceConfig
}

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
	apiVersion   string
	obj          runtime.Object
	oldObj       runtime.Object
}

// Controller object
type Controller struct {
	logger        *logrus.Entry
	clientset     kubernetes.Interface
	queue         workqueue.RateLimitingInterface
	informer      cache.SharedIndexInformer
	eventHandlers []handlers.Handler
}

func objName(obj interface{}) string {
	return reflect.TypeOf(obj).Name()
}
func init() {
	metric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "statemonitor_processed_changes_total",
		Help: "The total number of processed changes",
	},
		[]string{"Action", "Name", "Namespace", "Type"})
}

// TODO: we don't need the informer to be indexed
// Start prepares watchers and run their controllers, then waits for process termination signals
func Start(conf *config.Config, eventHandlers []handlers.Handler) {
	//TODO remove imput of evenhandlers and decide here
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	confDiff = conf.Diff
	namespaces = getNamespaces(kubeClient, &conf.NamespacesConfig)
	stopCh := make(chan struct{})
	ns := ""
	defer close(stopCh)

	if conf.Resource.CoreEvent.Enabled {
		allCoreEventsInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					options.FieldSelector = ""
					return kubeClient.CoreV1().Events(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					options.FieldSelector = ""
					return kubeClient.CoreV1().Events(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.Event{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, allCoreEventsInformer, objName(api_v1.Event{}), V1, conf.Resource.CoreEvent)

		go c.Run(stopCh)
	}

	if conf.Resource.Event.Enabled {

		allEventsInformer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					options.FieldSelector = ""
					return kubeClient.EventsV1().Events(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					options.FieldSelector = ""
					return kubeClient.EventsV1().Events(ns).Watch(context.Background(), options)
				},
			},
			&events_v1.Event{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, allEventsInformer, objName(events_v1.Event{}), EVENTS_V1, conf.Resource.Event)

		go c.Run(stopCh)
	}

	if conf.Resource.Pod.Enabled {
		pods := kubeClient.CoreV1().Pods(ns)
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					ll, err := pods.List(context.Background(), options)
					return ll, err
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					ww, err := pods.Watch(context.Background(), options)
					return ww, err
				},
			},
			&api_v1.Pod{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.Pod{}), V1, conf.Resource.Pod)

		go c.Run(stopCh)
	}

	if conf.Resource.HPA.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AutoscalingV1().HorizontalPodAutoscalers(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AutoscalingV1().HorizontalPodAutoscalers(ns).Watch(context.Background(), options)
				},
			},
			&autoscaling_v1.HorizontalPodAutoscaler{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(autoscaling_v1.HorizontalPodAutoscaler{}), AUTOSCALING_V1, conf.Resource.HPA)

		go c.Run(stopCh)

	}

	if conf.Resource.DaemonSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().DaemonSets(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().DaemonSets(ns).Watch(context.Background(), options)
				},
			},
			&apps_v1.DaemonSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(apps_v1.DaemonSet{}), APPS_V1, conf.Resource.DaemonSet)

		go c.Run(stopCh)
	}

	if conf.Resource.StatefulSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().StatefulSets(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().StatefulSets(ns).Watch(context.Background(), options)
				},
			},
			&apps_v1.StatefulSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(apps_v1.StatefulSet{}), APPS_V1, conf.Resource.StatefulSet)
		go c.Run(stopCh)
	}

	if conf.Resource.ReplicaSet.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().ReplicaSets(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().ReplicaSets(ns).Watch(context.Background(), options)
				},
			},
			&apps_v1.ReplicaSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(apps_v1.ReplicaSet{}), APPS_V1, conf.Resource.ReplicaSet)

		go c.Run(stopCh)
	}

	if conf.Resource.Services.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Services(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Services(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.Service{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.Service{}), V1, conf.Resource.Services)

		go c.Run(stopCh)
	}

	if conf.Resource.Deployment.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.AppsV1().Deployments(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.AppsV1().Deployments(ns).Watch(context.Background(), options)
				},
			},
			&apps_v1.Deployment{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(apps_v1.Deployment{}), APPS_V1, conf.Resource.Deployment)

		go c.Run(stopCh)
	}

	if conf.Resource.Namespace.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Namespaces().List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Namespaces().Watch(context.Background(), options)
				},
			},
			&api_v1.Namespace{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.Namespace{}), V1, conf.Resource.Namespace)

		go c.Run(stopCh)
	}

	if conf.Resource.ReplicationController.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ReplicationControllers(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ReplicationControllers(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.ReplicationController{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.ReplicationController{}), V1, conf.Resource.ReplicationController)

		go c.Run(stopCh)
	}

	if conf.Resource.Job.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.BatchV1().Jobs(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.BatchV1().Jobs(ns).Watch(context.Background(), options)
				},
			},
			&batch_v1.Job{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(batch_v1.Job{}), BATCH_V1, conf.Resource.Job)

		go c.Run(stopCh)
	}

	if conf.Resource.Node.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Nodes().List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Nodes().Watch(context.Background(), options)
				},
			},
			&api_v1.Node{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.Node{}), V1, conf.Resource.Node)

		go c.Run(stopCh)
	}

	if conf.Resource.ServiceAccount.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ServiceAccounts(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ServiceAccounts(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.ServiceAccount{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.ServiceAccount{}), V1, conf.Resource.ServiceAccount)

		go c.Run(stopCh)
	}

	if conf.Resource.ClusterRole.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.RbacV1().ClusterRoles().List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.RbacV1().ClusterRoles().Watch(context.Background(), options)
				},
			},
			&rbac_v1.ClusterRole{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(rbac_v1.ClusterRole{}), RBAC_V1, conf.Resource.ClusterRole)

		go c.Run(stopCh)
	}

	if conf.Resource.ClusterRoleBinding.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.RbacV1().ClusterRoleBindings().List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.RbacV1().ClusterRoleBindings().Watch(context.Background(), options)
				},
			},
			&rbac_v1.ClusterRoleBinding{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(rbac_v1.ClusterRoleBinding{}), RBAC_V1, conf.Resource.ClusterRoleBinding)

		go c.Run(stopCh)
	}

	if conf.Resource.PersistentVolume.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().PersistentVolumes().List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().PersistentVolumes().Watch(context.Background(), options)
				},
			},
			&api_v1.PersistentVolume{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.PersistentVolume{}), V1, conf.Resource.PersistentVolume)

		go c.Run(stopCh)
	}

	if conf.Resource.Secret.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().Secrets(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().Secrets(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.Secret{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.Secret{}), V1, conf.Resource.Secret)

		go c.Run(stopCh)
	}

	if conf.Resource.ConfigMap.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.CoreV1().ConfigMaps(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.CoreV1().ConfigMaps(ns).Watch(context.Background(), options)
				},
			},
			&api_v1.ConfigMap{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(api_v1.ConfigMap{}), V1, conf.Resource.ConfigMap)

		go c.Run(stopCh)
	}

	if conf.Resource.Ingress.Enabled {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
					return kubeClient.NetworkingV1().Ingresses(ns).List(context.Background(), options)
				},
				WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
					return kubeClient.NetworkingV1().Ingresses(ns).Watch(context.Background(), options)
				},
			},
			&networking_v1.Ingress{},
			0, //Skip resync
			cache.Indexers{},
		)

		c := newResourceController(kubeClient, eventHandlers, informer, objName(networking_v1.Ingress{}), NETWORKING_V1, conf.Resource.Ingress)

		go c.Run(stopCh)
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

// TODO: proper implementation of this function without the hack of multi ns
func newResourceController(client kubernetes.Interface, eventHandlers []handlers.Handler, informer cache.SharedIndexInformer, resourceType string, apiVersion string, resourceConfig config.ResourceConfig) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var eventWrapper EventWrapper
	eventWrapper.ResourceConfig = &resourceConfig

	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if resourceConfig.Enabled && (len(resourceConfig.IncludeEvenTypes) == 0 || slices.Contains(resourceConfig.IncludeEvenTypes, "add")) {
				var ok bool
				newEvent.namespace = "" // namespace retrived in processItem incase namespace value is empty
				newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
				newEvent.eventType = "create"
				newEvent.resourceType = resourceType
				newEvent.apiVersion = apiVersion
				newEvent.obj, ok = obj.(runtime.Object)
				if !ok {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot convert to runtime.Object for add on %v", obj)
				}
				if err != nil {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot get key for add on %v", obj)
					return
				}

				if !slices.Contains(namespaces, strings.Split(newEvent.key, "/")[0]) {
					logrus.Debugf("Skipping adding (namespaceconfig.ignore contains it) %v for %s", resourceType, newEvent.key)
					return
				}

				logrus.WithField("pkg", "statemonitor-"+resourceType).Infof("Processing add to %v: %s", resourceType, newEvent.key)

				eventWrapper.Event = newEvent
				queue.Add(eventWrapper)
			} else {
				logrus.Debugf("Skipping ADD (resource not enabled) %v for %s and is enabled - %t", resourceType, newEvent.key, resourceConfig.Enabled)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			if resourceConfig.Enabled && (len(resourceConfig.IncludeEvenTypes) == 0 || slices.Contains(resourceConfig.IncludeEvenTypes, "update")) {
				var ok bool
				newEvent.namespace = "" // namespace retrived in processItem incase namespace value is empty
				newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
				newEvent.eventType = "update"
				newEvent.resourceType = resourceType
				newEvent.apiVersion = apiVersion
				newEvent.obj, ok = new.(runtime.Object)
				if !ok {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot convert to runtime.Object for update on %v", new)
				}
				newEvent.oldObj, ok = old.(runtime.Object)
				if !ok {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot convert old to runtime.Object for update on %v", old)
				}

				if err != nil {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot get key for update on %v", old)
					return
				}

				if !slices.Contains(namespaces, strings.Split(newEvent.key, "/")[0]) {
					logrus.Debugf("Skipping updating(namespaceconfig.ignore contains it) %v for %s", resourceType, newEvent.key)
					return
				}

				logrus.WithField("pkg", "statemonitor-"+resourceType).Infof("Processing update to %v: %s", resourceType, newEvent.key)

				eventWrapper.Event = newEvent
				queue.Add(eventWrapper)
			} else {
				logrus.Debugf("Skipping UPDATE (resource not enabled) %v for %s and is enabled - %t", resourceType, newEvent.key, resourceConfig.Enabled)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if resourceConfig.Enabled && (len(resourceConfig.IncludeEvenTypes) == 0 || slices.Contains(resourceConfig.IncludeEvenTypes, "delete")) {
				var ok bool
				newEvent.namespace = "" // namespace retrived in processItem incase namespace value is empty
				newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				newEvent.eventType = "delete"
				newEvent.resourceType = resourceType
				newEvent.apiVersion = apiVersion
				newEvent.obj, ok = obj.(runtime.Object)
				if !ok {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot convert to runtime.Object for delete on %v", obj)
				}

				if err != nil {
					logrus.WithField("pkg", "statemonitor-"+resourceType).Errorf("cannot get key for delete on %v", obj)
					return
				}

				if !slices.Contains(namespaces, strings.Split(newEvent.key, "/")[0]) {
					logrus.Debugf("Skipping deletion (namespaceconfig.ignore contains it) %v for %s", resourceType, newEvent.key)
					return
				}

				logrus.WithField("pkg", "statemonitor-"+resourceType).Infof("Processing delete to %v: %s", resourceType, newEvent.key)
				eventWrapper.Event = newEvent
				queue.Add(eventWrapper)
			} else {
				logrus.Debugf("Skipping deletion (resource not enabled) %v for %s and is enabled - %t", resourceType, newEvent.key, resourceConfig.Enabled)
			}
		},
	})

	return &Controller{
		logger:        logrus.WithField("pkg", resourceType+"-statemonitor"),
		clientset:     client,
		informer:      informer,
		queue:         queue,
		eventHandlers: eventHandlers,
	}
}

// Run starts the statemonitor controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Starting statemonitor controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	c.logger.Info("statemonitor controller synced and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(EventWrapper))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

/* TODOs
- Enhance event creation using client-side cacheing machanisms - pending
- Enhance the processItem to classify events - done
- Send alerts correspoding to events - done
*/

func (c *Controller) processItem(eventWrapper EventWrapper) error {
	// NOTE that obj will be nil on deletes!
	newEvent := eventWrapper.Event
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)

	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", newEvent.key, err)
	}
	// get object's metedata
	objectMeta := utils.GetObjectMetaData(obj)

	// hold status type for default critical alerts
	var status string

	// namespace retrived from event key incase namespace value is empty
	if newEvent.namespace == "" && strings.Contains(newEvent.key, "/") {
		substring := strings.Split(newEvent.key, "/")
		newEvent.namespace = substring[0]
		newEvent.key = substring[1]
	} else {
		newEvent.namespace = objectMeta.Namespace
	}

	// process events based on its type
	switch newEvent.eventType {
	case "create":
		// compare CreationTimestamp and serverStartTime and alert only on latest events
		// Could be Replaced by using Delta or DeltaFIFO
		if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			switch newEvent.resourceType {
			case "NodeNotReady":
				status = "Danger"
			case "NodeReady":
				status = "Normal"
			case "NodeRebooted":
				status = "Danger"
			case "Backoff":
				status = "Danger"
			default:
				status = "Normal"
			}
			kbEvent := event.StatemonitorEvent{
				Name:       newEvent.key,
				Namespace:  newEvent.namespace,
				Kind:       newEvent.resourceType,
				ApiVersion: newEvent.apiVersion,
				Status:     status,
				Reason:     "Created",
			}

			for _, handler := range c.eventHandlers {
				handler.Handle(kbEvent)
			}

			handleMetric(newEvent)
			return nil
		}
	case "update":
		switch newEvent.resourceType {
		case "Backoff":
			status = "Danger"
		default:
			status = "Warning"
		}

		kbEvent := event.StatemonitorEvent{
			Name:       newEvent.key,
			Namespace:  newEvent.namespace,
			Kind:       newEvent.resourceType,
			ApiVersion: newEvent.apiVersion,
			Status:     status,
			Reason:     "Updated",
			Diff:       compareObjects(eventWrapper),
		}

		if kbEvent.Diff == "" {
			logrus.Printf("No diff( or ingored paths) found for %s", newEvent.key)
			//skipping metrics here as there is no valuable diff
			return nil
		}

		for _, handler := range c.eventHandlers {
			handler.Handle(kbEvent)
		}
		handleMetric(newEvent)
		return nil
	case "delete":
		kbEvent := event.StatemonitorEvent{
			Name:       newEvent.key,
			Namespace:  newEvent.namespace,
			Kind:       newEvent.resourceType,
			ApiVersion: newEvent.apiVersion,
			Status:     "Danger",
			Reason:     "Deleted",
		}

		for _, handler := range c.eventHandlers {
			handler.Handle(kbEvent)
		}
		handleMetric(newEvent)
		return nil
	}
	return nil
}

// compareObjects compares two objects and returns the diff
func compareObjects(ew EventWrapper) string {
	var patch jsondiff.Patch
	var err error
	ignorePath := append(confDiff.IgnorePath, ew.ResourceConfig.IgnorePath...)
	e := ew.Event
	oldObjj := e.oldObj
	objj := e.obj

	if e.resourceType == "ConfigMap" {
		patch, err = compareConfigMaps(oldObjj, objj)
	}

	if patch == nil || err != nil {
		patch, err = jsondiff.Compare(oldObjj, objj, jsondiff.Ignores(ignorePath...))
	}

	//jsondiff.CompareJSON(source, target)
	if err != nil {
		logrus.Printf("Error in comparing objects %s", err)
	}
	b, err := json.MarshalIndent(patch, "", "    ")
	if err != nil {
		logrus.Printf("Error in marshalling patch %s", err)
	}
	if b == nil || string(b) == "null" {
		return ""
	}
	return string(b)
}

func compareConfigMaps(old runtime.Object, new runtime.Object) (jsondiff.Patch, error) {

	//Dynamic extraction of data from configmap
	keys := make([]string, 0)
	for k, _ := range old.(*api_v1.ConfigMap).Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	if len(keys) == 0 {
		return nil, fmt.Errorf("error in extracting data from configmap")
	}
	k := keys[0]
	if !strings.Contains(k, ".json") {
		return nil, fmt.Errorf("error in extracting data from configmap")
	}

	oldData, oldSuccess := old.(*api_v1.ConfigMap).Data[k]
	newData, newSuccess := new.(*api_v1.ConfigMap).Data[k]
	if !oldSuccess || !newSuccess {
		return nil, fmt.Errorf("error in extracting data from configmap")
	}

	oldDataStr := strings.ReplaceAll(oldData, "\\", "")
	newDataStr := strings.ReplaceAll(newData, "\\", "")

	return jsondiff.CompareJSON([]byte(oldDataStr), []byte(newDataStr))
}

// getNamespaces returns the namespaces to watch based on the configiration provided *NamespacesConfig
func getNamespaces(clientset kubernetes.Interface, namespacesConfig *config.NamespacesConfig) []string {

	if namespacesConfig != nil && len(namespacesConfig.Include) > 0 {
		return namespacesConfig.Include
	}

	//Get all namespaces
	var namespaces []string
	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), meta_v1.ListOptions{})
	if err != nil {
		logrus.Errorf("Error in getting namespaces %s", err)
	}

	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	//Exclude namespaces from all namespaces
	if namespacesConfig != nil && len(namespacesConfig.Exclude) > 0 {
		for _, ns := range namespacesConfig.Exclude {
			for i, n := range namespaces {
				if ns == n {
					logrus.Infof("Removing namespace %s from watchlist", ns)
					namespaces[i] = namespaces[len(namespaces)-1]
					namespaces = namespaces[:len(namespaces)-1]
				}
			}
		}
	}

	logrus.Infof("Namespaces to watch %v", namespaces)
	return namespaces
}
func handleMetric(newEvent Event) {
	mu.Lock()
	defer mu.Unlock()
	metric.WithLabelValues([]string{newEvent.eventType, newEvent.key, newEvent.namespace, newEvent.resourceType}...).Inc()
}
