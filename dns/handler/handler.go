package handler

type Handler interface {
	SetStealthDNS() (bool, error)
	RemoveStealthDNS()
	GetUpstreamDNS() string
}
