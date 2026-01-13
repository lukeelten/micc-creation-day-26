package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_114362883")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "select1384045349",
			"maxSelect": 1,
			"name": "task",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"download",
				"convert",
				"process",
				"upload"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_114362883")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "select1384045349",
			"maxSelect": 1,
			"name": "task",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"download",
				"convert",
				"process",
				"upload"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
