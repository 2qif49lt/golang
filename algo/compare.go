package algo

type Compare interface {
	Less(r interface{}) bool
	Equal(r interface{}) bool
	More(r interface{}) bool
}
type Lesser interface {
	Less(r interface{}) bool
}
type Equaler interface {
	Equal(r interface{}) bool
}
type Morer interface {
	More(r interface{}) bool
}
