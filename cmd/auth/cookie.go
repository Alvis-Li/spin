package auth

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "sigs.k8s.io/yaml"

    "github.com/spinnaker/spin/config"
    "github.com/spinnaker/spin/config/auth"
)

type cookieOptions struct {
    *authOptions
}

var (
    cookieShort   = "Store a session cookie in the spin configuration"
    cookieLong    = "Store a raw Cookie header value (e.g. 'SESSION=abcd1234') in ~/.spin/config so that all subsequent spin commands include it automatically."
    cookieExample = "spin auth cookie \"SESSION=abcd1234\""
)

// NewCookieCmd creates a new `auth cookie` command.
func NewCookieCmd(options *authOptions) *cobra.Command {
    opts := &cookieOptions{authOptions: options}

    cmd := &cobra.Command{
        Use:     "cookie <cookie-string>",
        Short:   cookieShort,
        Long:    cookieLong,
        Example: cookieExample,
        Args:    cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return setCookie(opts, args[0])
        },
    }

    return cmd
}

func setCookie(options *cookieOptions, cookie string) error {
    // Determine config file location. Respect --config flag if provided.
    var configPath string

    // Fallback to default location ($HOME/.spin/config).
    {
        home, err := os.UserHomeDir()
        if err != nil {
            return fmt.Errorf("unable to determine user home directory: %w", err)
        }
        configPath = filepath.Join(home, ".spin", "config")
    }


    // Ensure the parent directory exists.
    if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
        return fmt.Errorf("unable to create config directory: %w", err)
    }

    // Load existing configuration if present.
    var cfg config.Config
    if data, err := ioutil.ReadFile(configPath); err == nil {
        if err := yaml.UnmarshalStrict([]byte(os.ExpandEnv(string(data))), &cfg); err != nil {
            return fmt.Errorf("failed to parse existing config: %w", err)
        }
    }

    if cfg.Auth == nil {
        cfg.Auth = &auth.Config{}
    }
    cfg.Auth.Cookie = cookie

    // Marshal back to YAML.
    buf, err := yaml.Marshal(&cfg)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    // Preserve existing file mode or default to 0600.
    mode := os.FileMode(0600)
    if info, err := os.Stat(configPath); err == nil {
        mode = info.Mode()
    }

    if err := ioutil.WriteFile(configPath, buf, mode); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    if options.Ui != nil {
        options.Ui.Info(fmt.Sprintf("Cookie updated and stored in %s", configPath))
    }

    return nil
} 