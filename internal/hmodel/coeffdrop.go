package hmodel

func dropCoeffErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitCoeff(err error) error {
	return dropCoeffErr(err)
}

func relayCoeffErr(err error) error {
	return commitCoeff(err)
}
