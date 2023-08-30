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
		msg = createBoxlikeOutput(e)
	}
	return msg
}

func createBoxlikeOutput(e *DiffWatchEvent) string {
	var sb strings.Builder
	sb.Grow(1200)

	const col1Width = 12
	var col2Width = 15

	if len(e.Name) > col2Width {
		col2Width = len(e.Name) + 2
	}
	sb.WriteString(e.Diff + "\n")

	dataRow(&sb, col1Width, col2Width, "Type", e.Kind)
	dataRow(&sb, col1Width, col2Width, "Name", e.Name)
	dataRow(&sb, col1Width, col2Width, "Action", e.Reason)
	dataRow(&sb, col1Width, col2Width, "Namespace", e.Namespace)
	dataRow(&sb, col1Width, col2Width, "Status", e.Status)
	sb.WriteString(fmt.Sprintf("+%s+%s+\n", strings.Repeat("-", col1Width), strings.Repeat("-", col2Width)))

	return sb.String()
}

// Write the data rows
func dataRow(sb *strings.Builder, col1Width int, col2Width int, description string, value string) {

	sb.WriteString(fmt.Sprintf("+%s+%s+\n",
		strings.Repeat("-", col1Width), strings.Repeat("-", col2Width)))

	sb.WriteString(fmt.Sprintf("| %-"+fmt.Sprintf("%d", col1Width-2)+"s | %-"+fmt.Sprintf("%d", col2Width-2)+"s |\n", description, value))
}
