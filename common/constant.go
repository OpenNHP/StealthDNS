package common

const (
	StealthDnsIp        = "127.0.0.1"
	DnsUdpPort          = 53
	NhpDomainNameSuffix = ".nhp"
)

const (
	Type_A = iota
	Type_AAAA
)

const (
	AgentInit          = "nhp_agent_init"
	AgentClose         = "nhp_agent_close"
	AgentFreeCString   = "nhp_free_cstring"
	AgentKnockResource = "nhp_agent_knock_resource"
)
