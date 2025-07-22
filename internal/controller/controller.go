package controller

type Harbor struct {
	Enable   bool   `yaml:"enable" json:"enable"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Scheme   string `json:"scheme" yaml:"scheme"`
}
type Prometheus struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Host   string `json:"host" yaml:"host"`
	Scheme string `json:"scheme" yaml:"scheme"`
}
type System struct {
	Addr        string     `json:"addr" yaml:"addr"`
	Provisioner string     `json:"provisioner" yaml:"provisioner"`
	Harbor      Harbor     `json:"harbor" yaml:"harbor"`
	Prometheus  Prometheus `json:"prometheus" yaml:"prometheus"`
}

type Server struct {
	System System `json:"system" yaml:"system"`
}
