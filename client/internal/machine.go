package internal

import (
	"github.com/flowshield/flowshield/client/internal/config"
	"github.com/flowshield/flowshield/client/pkg/util/uuid"
)

// InitMachine initialize the machine id
func InitMachine() error {
	machine := config.C.Machine
	mac, err := machine.Read()
	if err != nil {
		machine.SetMachineId(uuid.MustString())
		err = machine.Write()
		if err != nil {
			return err
		}
	} else {
		machine.SetMachineId(mac.MachineId)
		machine.SetCookie(mac.Cookie)
	}
	return nil
}
