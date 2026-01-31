package entities

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     string
}
