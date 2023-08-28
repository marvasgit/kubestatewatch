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

package event

import (
	"fmt"
	"strings"
)

// DiffWatchEvent represent an event got from k8s api server
// Events from different endpoints need to be casted to DiffWatchEvent
// before being able to be handled by handler
type DiffWatchEvent struct {
	Namespace  string
	Kind       string
	ApiVersion string
	Component  string
	Host       string
	Reason     string
	Status     string
	Name       string
	Diff       string
}

// Message returns event message in standard format.
// included as a part of event packege to enhance code resuablity across handlers.
func (e *DiffWatchEvent) Message() (msg string) {
	// using switch over if..else, since the format could vary based on the kind of the object in future.
	switch e.Kind {
	case "namespace":
		msg = fmt.Sprintf(
			"A namespace `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "node":
		msg = fmt.Sprintf(
			"A node `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "cluster role":
		msg = fmt.Sprintf(
			"A cluster role `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "NodeReady":
		msg = fmt.Sprintf(
			"Node `%s` is Ready : \nNodeReady",
			e.Name,
		)
	case "NodeNotReady":
		msg = fmt.Sprintf(
			"Node `%s` is Not Ready : \nNodeNotReady",
			e.Name,
		)
	case "NodeRebooted":
		msg = fmt.Sprintf(
			"Node `%s` Rebooted : \nNodeRebooted",
			e.Name,
		)
	case "Backoff":
		msg = fmt.Sprintf(
			"Pod `%s` in `%s` Crashed : \nCrashLoopBackOff %s",
			e.Name,
			e.Namespace,
			e.Reason,
		)
	default:
		msg = formatDefaultMessage(e)
	}
	return msg
}

func formatDefaultMessage(e *DiffWatchEvent) string {
	maxLen := getlongerString(e.Name, e.Namespace)
	totalLen := maxLen + (40 - maxLen) + 17 //dont like magic numbers but this is all the spaces and chars in the message + Namespace str
	var sb strings.Builder
	sb.WriteString(e.Diff + "\n")
	sb.WriteString(strings.Repeat("-", totalLen) + "\n")
	sb.WriteString(strings.Repeat("-", totalLen) + "\n")
	sb.WriteString(fmt.Sprintf("| %-10s | %-40s |\n", "Type", e.Kind))
	sb.WriteString(fmt.Sprintf("| %-10s | %-40s |\n", "Name", e.Name))
	sb.WriteString(fmt.Sprintf("| %-10s | %-40s |\n", "Action", e.Reason))
	sb.WriteString(fmt.Sprintf("| %-10s | %-40s |\n", "Namespace", e.Namespace))
	sb.WriteString(fmt.Sprintf("| %-10s | %-40s |\n", "Status", e.Status))
	sb.WriteString(strings.Repeat("-", totalLen) + "\n")
	sb.WriteString(strings.Repeat("-", totalLen) + "\n")
	return sb.String()
}

func getlongerString(a, b string) int {
	if len(a) > len(b) {
		return len(a)
	}
	return len(b)
}
