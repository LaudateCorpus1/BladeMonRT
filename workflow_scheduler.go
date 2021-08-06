package main

import (
	"encoding/json"
	"github.com/microsoft/BladeMonRT/workflows"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
	winEvents "github.com/microsoft/BladeMonRT/windows_events"
)

/** Class for scheduling workflows. */
type WorkflowScheduler struct {
	schedules []interface{}
	logger *log.Logger
	subscriber winEvents.EventSubscriber
}

/** Class for the schedule description in the JSON. */
type ScheduleDescription struct {
    Name string `json:"name"`
    ScheduleType string `json:"schedule_type"`
    Workflow string `json:"workflow"`
    Enable bool `json:"enable"`

	// Attributes specific to events of type 'on_win_event'.
	WinEventSubscribeQueries [][]string `json:"win_event_subscribe_queries"`
}

/** Class that represents a query for subscribing to a windows event. */
type WinEventSubscribeQuery struct {
	channel string
	query string
}

func (workflowScheduler *WorkflowScheduler) addWinEventBasedSchedule(workflow workflows.InterfaceWorkflow, eventQueries []WinEventSubscribeQuery) {
	workflowScheduler.logger.Println("Workflow:", workflow)

	// Subscribe to the events that match the event queries specified.
	for _, eventQuery := range eventQueries {
		var eventSubscription *winEvents.EventSubscription = &winEvents.EventSubscription{
			Channel:        eventQuery.channel,
			Query:          eventQuery.query,
			SubscribeMethod: winEvents.EVT_SUBSCRIBE_TO_FUTURE_EVENTS,
			Callback:        workflowScheduler.subscriber.SubscriptionCallback,
			Context:         winEvents.CallbackContext{Workflow : workflow},
		}
		workflowScheduler.subscriber.CreateSubscription(eventSubscription)
	}
}

func newWorkflowScheduler(schedulesJson []byte, workflowFactory WorkflowFactory) *WorkflowScheduler {
	var subscriber winEvents.EventSubscriber = winEvents.NewEventSubscriber()
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("WorkflowScheduler")
	var workflowScheduler *WorkflowScheduler = &WorkflowScheduler{subscriber : subscriber, logger: logger}

	// Parse the schedules JSON and add the schedules to the workflow scheduler.
	var schedules map[string][]ScheduleDescription
	json.Unmarshal([]byte(schedulesJson), &schedules)
	for _, schedule := range schedules["schedules"] {
		switch schedule.ScheduleType {
			case "on_win_event":
				var workflow workflows.InterfaceWorkflow = workflowFactory.constructWorkflow(schedule.Workflow)	
				var eventQueries []WinEventSubscribeQuery = parseEventSubscribeQueries(schedule.WinEventSubscribeQueries)			
				workflowScheduler.addWinEventBasedSchedule(workflow, eventQueries) 
			default:
				workflowScheduler.logger.Println("Given schedule type not supported.")
		}
	}
	return workflowScheduler
}

func parseEventSubscribeQueries(eventQueries [][]string) []WinEventSubscribeQuery {
	var parsedEventQueries []WinEventSubscribeQuery
	// Parse each of the event queries into the 'WinEventSubscribeQuery' type.
	for _, eventQuery := range eventQueries {
		var parsedEventQuery = eventQuery
		var channel string = parsedEventQuery[0]
		var query string = parsedEventQuery[1]
		parsedEventQueries = append(parsedEventQueries, WinEventSubscribeQuery{channel : channel, query : query})
	}
	return parsedEventQueries
}