package ast

type Acceptor interface {
	Accept(visitor Visitor) error
}
