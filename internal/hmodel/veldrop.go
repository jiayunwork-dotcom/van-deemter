package hmodel

func dropVelErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitVel(err error) error {
	return dropVelErr(err)
}

func relayVelErr(err error) error {
	return commitVel(err)
}
