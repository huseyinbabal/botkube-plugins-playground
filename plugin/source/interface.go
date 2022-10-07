package plugin

type Source interface {
	Consume(ch chan interface{}) error
}
