package subscriptions

type Subscription struct {
	Id    int
	Email string
}

type IRepository interface {
	GetAll() ([]*Subscription, error)
	Add(email string) error
	Exists(email string) (bool, error)
}

type IService interface {
	Add(email string) (bool, error)
}
