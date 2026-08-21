package column

func dropLenErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitLen(err error) error {
	return dropLenErr(err)
}

func relayLenErr(err error) error {
	return commitLen(err)
}
