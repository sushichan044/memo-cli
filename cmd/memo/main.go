package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/sushichan044/memo-cli/internal/config"
	"github.com/sushichan044/memo-cli/internal/memo"
	"github.com/sushichan044/memo-cli/version"
)

type (
	CLIContext struct {
		cfg *config.Config
	}

	CLI struct {
		Version kong.VersionFlag `short:"v" help:"Show version."`

		New NewCmd `cmd:"new" help:"Create a new memo."`
		// List ListCmd `cmd:"list" help:"List all memos."` TODO: add go-fzf integration
	}

	NewCmd struct {
		Name string `arg:"" optional:"" help:"Memo name (default: timestamp HH-MM-SS)"`
	}
)

func (c *NewCmd) Run(ctx *CLIContext) error {
	creator := memo.New(ctx.cfg)

	// Check gitignore and print warning if needed
	if warning := creator.CheckGitignore(); warning != "" {
		fmt.Fprintln(os.Stderr, warning)
		fmt.Fprintln(os.Stderr) // blank line
	}

	path, err := creator.Create(c.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	// Output success message to stderr
	fmt.Fprintf(os.Stderr, "✅ Memo created at: %s\n", path)

	// Output path to stdout (for piping)
	fmt.Println(path) //nolint:forbidigo // stdout output is intentional for piping

	return nil
}

func main() {
	ctx := kong.Parse(&CLI{},
		kong.Vars{
			"version": fmt.Sprintf("memo-cli %s", version.Get()),
		},
		kong.Name("memo"),
		kong.Description("A CLI tool to create and manage markdown memos"),
		kong.UsageOnError(),
	)

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	ctx.FatalIfErrorf(ctx.Run(&CLIContext{cfg: cfg}))
}
