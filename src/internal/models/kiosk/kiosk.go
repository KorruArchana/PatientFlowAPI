package kiosk

type Kiosk struct {
	PK            string `json:"pk"`
	SK            string `json:"sk"`
	KioskName     string `json:"kioskName"`
	AdminPassword string `json:"adminPassword"`
	ScreenTimeout string `json:"screenTimeOut"`
	MachineID     string `json:"machineID"`
}
