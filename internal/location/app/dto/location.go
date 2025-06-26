package dto

type Country struct {
	Name string
	Code string
}

type State struct {
	Name string
	Code string
}

type Address struct {
	Line1 string
	Line2 string
	Line3 string
}

type Location struct {
	Id            int
	ZipCode       string
	Country       Country
	State         State
	City          string
	Address       Address
	IsResidential bool
	Name          string
	Phone         string
}
