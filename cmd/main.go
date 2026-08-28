package main

import (
	"os"

	"github.com/elliotwutingfeng/passwordsim"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

type results []result

func (u results) Len() int {
	return len(u)
}
func (u results) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
func (u results) Less(i, j int) bool {
	if u[i].score < u[j].score {
		return true
	}
	if u[i].score == u[j].score {
		return u[i].password < u[j].password
	}
	return false
}

type result struct {
	password string
	score    float64
}

func main() {

	var passwords string
	var output string
	var passwordToCheck string
	var threshold float64

	app := &cli.App{
		Name: "passwordsim",
		Description: `passwordsim lets you search for passwords similar to your specified password in any passwords dataset.
The similarity metric used is the Damerau-Levenshtein distance.

EXAMPLE (for UNIX systems): passwordsim -f 'test/passwords.txt' -o 'output.txt' -p 'correct horse battery staple' -t 0.3

EXPLANATION: The above command searches the file located at 'test/passwords.txt' for passwords with a normalised
Damerau-Levenshtein distance no more than 0.3 relative to the specified password 'correct horse battery staple'.

These passwords and their respective normalised Damerau-Levenshtein distance scores are saved to 'output.txt'.

Normalised Damerau-Levenshtein distance scores range from 0.0 to 1.0.

The smaller the normalised Damerau-Levenshtein distance between 2 passwords, the more similar they are to each other.
`,
		Usage: `Search any leaked passwords dataset for passwords similar to your specified password.
Dataset format: Text file, one password per line.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Value:       "",
				Usage:       "Path to `FILE` containing leaked passwords dataset.",
				Destination: &passwords,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "output.txt",
				Usage:       "Path to `FILE` where similar passwords and their normalised Damerau-Levenshtein distance scores are to be stored.",
				Destination: &output,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "`PASSWORD` to check.",
				Destination: &passwordToCheck,
			},
			&cli.Float64Flag{
				Name:    "threshold",
				Aliases: []string{"t"},
				Value:   0.3,
				Usage: "Maximum normalised Damerau-Levenshtein distance score `THRESHOLD`." +
					" Out-of-range value will be clamped within range [0.0, 1.0].",
				Destination: &threshold,
			},
		},
		Action: func(cCtx *cli.Context) error {
			var fileFlag bool
			var passwordFlag bool
			for _, flag := range cCtx.LocalFlagNames() {
				if flag == "help" {
					return nil
				}
				if flag == "file" {
					fileFlag = true
				}
				if flag == "password" {
					passwordFlag = true
				}
			}
			if cCtx.NArg() != 0 {
				color.HiRed("Input Error: Use flags instead of arguments. Check `--help` for more details.")
				os.Exit(1)
			}
			if !fileFlag {
				color.HiRed("Input Error: No dataset file specified. Is the `-f` flag missing?")
				os.Exit(1)
			}
			if !passwordFlag {
				color.HiRed("Input Error: No password to check specified. Is the `-p` flag missing?")
				os.Exit(1)
			}
			passwordsim.CheckPasswords(passwords, output, passwordToCheck, threshold)
			return nil
		},
	}
	app.Run(os.Args)
}
