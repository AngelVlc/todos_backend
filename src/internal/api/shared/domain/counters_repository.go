package domain

type CountersRepository interface {
	FindByName(name string) (*Counter, error)
	Create(counter *Counter) error
	Update(counter *Counter) error
}
