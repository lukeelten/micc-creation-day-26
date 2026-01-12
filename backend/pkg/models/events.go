package models

import "github.com/pocketbase/pocketbase/core"

type EventsTypeOptions string

const (
	EventsTypeInfo  EventsTypeOptions = "info"
	EventsTypeWarn  EventsTypeOptions = "warn"
	EventsTypeError EventsTypeOptions = "error"
)

type EventsRecord struct {
	AdditionalData any               `json:"additionalData,omitempty"`
	AdditionalText string            `json:"additionalText,omitempty"`
	Created        string            `json:"created"`
	Description    string            `json:"description,omitempty"`
	ID             string            `json:"id"`
	Run            string            `json:"run"`
	Title          string            `json:"title"`
	Type           EventsTypeOptions `json:"type"`
	Updated        string            `json:"updated"`
}

func CreateEvent(app core.App, event *EventsRecord) error {
	collection, err := app.FindCollectionByNameOrId(CollectionsEvents)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)

	record.Set("run", event.Run)
	record.Set("title", event.Title)
	record.Set("type", string(event.Type))
	record.Set("description", event.Description)
	record.Set("additionalText", event.AdditionalText)
	record.Set("additionalData", event.AdditionalData)

	return app.Save(record)
}
