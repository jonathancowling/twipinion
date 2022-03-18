package iampolicy

type Principal struct {
	AWS []string `json:",omitempty"`
	CanonicalUser []string `json:",omitempty"`
	Federated []string `json:",omitempty"`
	Service []string `json:",omitempty"`
}

type Statement struct {
	Principal *Principal `json:",omitempty"`
	NotPrincipal *Principal`json:",omitempty"`
	Effect string `json:",omitempty"`
	Action []string `json:",omitempty"`
	NotAction []string `json:",omitempty"`
	Resource []string `json:",omitempty"`
	NotResource []string `json:",omitempty"`
	Condition map[string]map[string]string `json:",omitempty"`
}

type IamPolicy struct {
	Version string
	Statement Statement
}
