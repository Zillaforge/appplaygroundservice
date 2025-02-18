package quiz

const (
	TypeNameString            = "string"
	TypeNamePassword          = "password"
	TypeNameBoolean           = "boolean"
	TypeNameInt               = "int"
	TypeNameEnum              = "enum"
	TypeNameArray             = "array"
	TypeNameVPSNetwork        = "vpsNetwork"
	TypeNameVPSFlavor         = "vpsFlavor"
	TypeNameVPSGPUFlavor      = "vpsGPUFlavor"
	TypeNameVPSvGPUFlavor     = "vpsvGPUFlavor"
	TypeNameVPSSecurityGroups = "vpsSecurityGroups"
	TypeNameVPSKeypair        = "vpsKeypair"
	TypeNameVPSVolume         = "vpsVolume"
	TypeNameVPSBootVolume     = "vpsBootVolume"
	TypeNamePort              = "port"
	TypeNameSSHPort           = "sshPort"
)

type Questions struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	Label       string       `json:"label"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Variable    string       `json:"variable"`
	Group       string       `json:"group"`
	Required    bool         `json:"required"`
	Order       int          `json:"order"`
	Default     *interface{} `json:"default"`
	Options     *[]string    `json:"options"`
}

type Answers struct {
	Answers []Answer `json:"answers"`
}

type Answer struct {
	Question
	RawValue    interface{} `json:"rawValues"`
	Value       interface{} `json:"values"`
	DisplayName interface{} `json:"displayName"`
}
