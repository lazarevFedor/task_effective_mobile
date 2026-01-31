// Package entities contains domain entities used across the application.
//
// Entities are plain Go structs that represent business objects persisted by
// repositories and transported via HTTP handlers.
package entities

// Subscription represents a user's subscription to a service.
//
// ID is the database identifier. ServiceName is the name of the subscribed
// service (for example "netflix"). Price is an integer amount (in cents or
// the minimal currency unit as used by your application). UserID references
// the owner of the subscription. StartDate and EndDate are formatted as
// "MM-YYYY" when exposed via the API; EndDate may be empty to indicate an
// open-ended subscription.
type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     string
}
