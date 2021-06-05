// +build windows

package passport

import "golang.org/x/sys/windows/registry"

// getMachineID returns the windows host's unique machine GUID.
// https://github.com/denisbrodbeck/machineid/blob/master/id_windows.go
func getMachineID() ([]byte, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
	if err != nil {
		return nil, err
	}
	defer k.Close()

	bytes, _, err := k.GetBinaryValue("MachineGuid")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
