package goengine

type HttpAction interface {
	Bind(* HttpRoute)
}