package app

import (
	"github.com/tepleton/tepleton-sdk/x/bank"
)

// initRoutes() happens after initCapKeys(), initStores(), and initSDKApp().
func (app *BasecoinApp) initRoutes() {
	var router = app.BaseApp.Router()
	var accStore = app.accStore

	// All handlers must be added here.
	// The order matters.
	router.AddRoute("bank", bank.NewHandler(accStore))
}
