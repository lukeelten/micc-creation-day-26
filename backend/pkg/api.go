package pkg

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func SetupApi(pb *pocketbase.PocketBase) error {

	pb.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Setup API routes here

		return e.Next()
	})

	return nil
}
