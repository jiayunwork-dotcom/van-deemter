package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"van-deemter/internal/evalx"
)

const usage = `van-deemter: van Deemter plate-height and plate-count accounting.

Reads a packed column case from a JSON file with fields a (eddy diffusion, m),
b (longitudinal diffusion, m^2/s), c (mass-transfer resistance, s), length (m)
and optionally n_req, kprime and alpha. For a mobile-phase velocity u (m/s)
computes the plate height H(u) = A + B/u + C*u, the theoretical plate count
N = L/H(u), and when n_req is given the maximum allowed H for that requirement
and whether the current u satisfies it. When kprime and alpha are both given,
computes the simplified resolution Rs = (sqrt(N)/4) * ((alpha-1)/alpha) *
(kprime/(1+kprime)); without kprime and alpha, Rs is never fabricated.

usage:
  van-deemter eval <input.json> --u <m/s>
  van-deemter eval <input.json> --scan
  van-deemter help

--u  prints H, N (and H_max/meets, Rs when the case provides them).
--scan prints u_opt = sqrt(B/C), H_min = A + 2*sqrt(B*C), a grid spot-check
  and the central derivative at u_opt.

boundaries: A, B, C, length and u must all be positive; u<=0, a negative
coefficient or a non-positive length are reported on stderr with a non-zero
exit code. n_req must be positive. kprime and alpha must be given together.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "eval":
		if err := runEval(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "van-deemter: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "van-deemter: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runEval(args []string) error {
	fs := flag.NewFlagSet("eval", flag.ContinueOnError)
	fs.SetOutput(flagDiscard{})
	u := fs.Float64("u", 0, "mobile-phase velocity u in m/s")
	scanMode := fs.Bool("scan", false, "print u_opt and H_min instead of H and N")
	flagArgs, files := splitFlags(args)
	if err := fs.Parse(flagArgs); err != nil {
		return err
	}
	if len(files) < 1 {
		return fmt.Errorf("eval needs a case file")
	}
	c, err := evalx.LoadCase(files[0])
	if err != nil {
		return err
	}
	if *scanMode {
		return evalx.RunScan(os.Stdout, c)
	}
	if *u <= 0 {
		return fmt.Errorf("eval needs --u <m/s> (positive) or --scan")
	}
	return evalx.RunEval(os.Stdout, c, *u)
}

func splitFlags(args []string) (flags, files []string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			files = append(files, args[i+1:]...)
			break
		}
		if strings.HasPrefix(a, "-") && a != "-" {
			flags = append(flags, a)
			if !strings.Contains(a, "=") {
				switch a {
				case "--u", "-u":
					if i+1 < len(args) {
						flags = append(flags, args[i+1])
						i++
					}
				}
			}
			continue
		}
		files = append(files, a)
	}
	return flags, files
}

type flagDiscard struct{}

func (flagDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
