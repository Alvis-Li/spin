package auth

import (
    "github.com/spf13/cobra"

    "github.com/spinnaker/spin/cmd"
)

// authOptions bundles common options needed by auth subcommands.
type authOptions struct {
    *cmd.RootOptions
}

var (
    authShort   = "Authentication related commands"
    authLong    = "Manage authentication configuration for spin."
    authExample = "spin auth cookie \"SESSION=abcd1234\""
)

// NewAuthCmd returns the root authentication command which hosts subcommands
// like `cookie`.
func NewAuthCmd(rootOptions *cmd.RootOptions) *cobra.Command {
    options := &authOptions{RootOptions: rootOptions}

    cmd := &cobra.Command{
        Use:     "auth",
        Short:   authShort,
        Long:    authLong,
        Example: authExample,
    }

    // Register subcommands.
    cmd.AddCommand(NewCookieCmd(options))

    return cmd
} 