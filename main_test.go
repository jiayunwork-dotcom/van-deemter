package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, run func() error) (string, string, error) {
	t.Helper()
	oldOut, oldErr := os.Stdout, os.Stderr
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}
	os.Stdout, os.Stderr = outW, errW
	runErr := run()
	outW.Close()
	errW.Close()
	os.Stdout, os.Stderr = oldOut, oldErr
	stdout, _ := io.ReadAll(outR)
	stderr, _ := io.ReadAll(errR)
	outR.Close()
	errR.Close()
	return string(stdout), string(stderr), runErr
}

func TestMainEvalPrintsHAndN(t *testing.T) {
	stdout, _, err := captureOutput(t, func() error {
		return runEval([]string{"example/packed.json", "--u", "0.02"})
	})
	if err != nil {
		t.Fatalf("runEval error: %v", err)
	}
	if !strings.Contains(stdout, "H ") {
		t.Errorf("stdout %q missing H line", stdout)
	}
	if !strings.Contains(stdout, "N ") {
		t.Errorf("stdout %q missing N line", stdout)
	}
	if !strings.Contains(stdout, "3.44") || !strings.Contains(stdout, "e-05") {
		t.Errorf("stdout %q missing expected H value 3.44e-05", stdout)
	}
	if !strings.Contains(stdout, "7267") {
		t.Errorf("stdout %q missing expected N value 7267", stdout)
	}
	if !strings.Contains(stdout, "meets") {
		t.Errorf("stdout %q missing requirement judgement for packed.json", stdout)
	}
}

func TestMainScanPrintsOptimum(t *testing.T) {
	stdout, _, err := captureOutput(t, func() error {
		return runEval([]string{"example/packed.json", "--scan"})
	})
	if err != nil {
		t.Fatalf("runEval --scan error: %v", err)
	}
	if !strings.Contains(stdout, "u_opt") {
		t.Errorf("stdout %q missing u_opt line", stdout)
	}
	if !strings.Contains(stdout, "H_min") {
		t.Errorf("stdout %q missing H_min line", stdout)
	}
	if !strings.Contains(stdout, "0.005774") {
		t.Errorf("stdout %q missing expected u_opt 0.005774", stdout)
	}
	if !strings.Contains(stdout, "2.9543e-05") && !strings.Contains(stdout, "2.954e-05") {
		t.Errorf("stdout %q missing expected H_min 2.954e-05", stdout)
	}
}

func TestMainInvalidVelocityExitsNonZero(t *testing.T) {
	stdout, _, err := captureOutput(t, func() error {
		return runEval([]string{"example/packed.json"})
	})
	if err == nil {
		t.Error("eval without --u and without --scan expected error, got nil")
	} else if !strings.Contains(err.Error(), "eval needs") {
		t.Errorf("error %v must explain the missing flag", err)
	}
	stdout2, _, err2 := captureOutput(t, func() error {
		return runEval([]string{"example/packed.json", "--u", "0"})
	})
	if err2 == nil {
		t.Error("eval --u 0 expected error, got nil")
	} else if !strings.Contains(err2.Error(), "eval needs") {
		t.Errorf("error %v must reject a non-positive velocity", err2)
	}
	if strings.TrimSpace(stdout) != "" || strings.TrimSpace(stdout2) != "" {
		t.Errorf("invalid input must not print results (got %q / %q)", stdout, stdout2)
	}
	_, _, err3 := captureOutput(t, func() error {
		return runEval([]string{"example/packed.json", "--u", "-0.1"})
	})
	if err3 == nil {
		t.Error("eval --u -0.1 expected error, got nil")
	}
}
