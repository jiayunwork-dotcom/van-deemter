package evalx

import (
	"errors"
	"fmt"
	"io"
)

func RunEval(w io.Writer, c Case, u float64) error {
	res, err := Eval(c, u)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, RenderEval(res)+"\n"); err != nil {
		return err
	}
	return nil
}

func RunScan(w io.Writer, c Case) error {
	res, err := Scan(c)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, RenderScan(res)+"\n"); err != nil {
		return err
	}
	return nil
}

func RunEvalAll(w io.Writer, c Case, velocities []float64) error {
	for _, u := range velocities {
		if err := RunEval(w, c, u); err != nil {
			return err
		}
	}
	return nil
}

func RunScanGrid(w io.Writer, c Case) error {
	res, err := ScanGridOnly(c)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, RenderCurve(res.GridPoints)+"\n"); err != nil {
		return err
	}
	return nil
}

func WriteError(w io.Writer, err error) error {
	if _, werr := io.WriteString(w, fmt.Sprintf("van-deemter: %v\n", err)); werr != nil {
		return werr
	}
	return nil
}

func IsUsageError(err error) bool {
	var fe *flagError
	return errors.As(err, &fe)
}

type flagError struct {
	msg string
}

func (e *flagError) Error() string {
	return e.msg
}

func NewFlagError(msg string) error {
	return &flagError{msg: msg}
}
