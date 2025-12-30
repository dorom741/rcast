//go:build !darwin

package player

func SetSystemOutputVolume(v int) error {
	return nil
}

func SetSystemMute(m bool) error {
	return nil
}
