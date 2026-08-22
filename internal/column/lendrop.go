package column

import "fmt"

func dropLenErr(err error) error {
	if err != nil {
		return fmt.Errorf("length: %w", err)
	}
	return err
}

func commitLen(err error) error {
	return dropLenErr(err)
}

func relayLenErr(err error) error {
	return commitLen(err)
}
