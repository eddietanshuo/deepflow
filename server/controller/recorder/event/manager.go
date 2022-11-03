/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"time"

	"github.com/deepflowys/deepflow/server/controller/recorder/cache"
	"github.com/deepflowys/deepflow/server/controller/recorder/constraint"
	"github.com/deepflowys/deepflow/server/libs/eventapi"
	"github.com/deepflowys/deepflow/server/libs/queue"
)

type EventManager[CT constraint.CloudModel, MT constraint.MySQLModel, BT constraint.DiffBase[MT]] struct {
	resourceType string
	ToolDataSet  cache.ToolDataSet
	Queue        *queue.OverwriteQueue
	EventProducer[CT, MT, BT]
}

func (e *EventManager[CT, MT, BT]) createAndPutEvent(eventType string, resourceType, resourceID int, resourceName, description string) {
	event := e.create()
	event.Time = time.Now().Unix()
	event.Type = eventType
	event.ResourceType = uint32(resourceType)
	event.ResourceID = uint32(resourceID)
	event.ResourceName = resourceName
	event.Description = description
	e.put(event)
}

func (e *EventManager[CT, MT, BT]) create() *eventapi.ResourceEvent {
	return eventapi.AcquireResourceEvent()
}

func (e *EventManager[CT, MT, BT]) put(event *eventapi.ResourceEvent) {
	err := e.Queue.Put(event)
	if err != nil {
		log.Error(putEventIntoQueueFailed(e.resourceType, err))
	}
	log.Infof("put %s event: %+v into shared queue", e.resourceType, event)
}

type EventProducer[CT constraint.CloudModel, MT constraint.MySQLModel, BT constraint.DiffBase[MT]] interface {
	ProduceByAdd([]*MT)
	ProduceByUpdate(*CT, BT)
	ProduceByDelete([]string)
}