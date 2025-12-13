package app

import (
	"github.com/codoworks/sshto/internal/config"
	"github.com/codoworks/sshto/internal/ssh"
)

// App orchestrates the sshto application
type App struct {
	Config    *config.Config
	SSHClient *ssh.Client
}

// New creates a new App instance
func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:    cfg,
		SSHClient: ssh.NewClient(),
	}, nil
}

// Connect establishes an SSH connection to the named server with optional overrides
func (a *App) Connect(serverName string, opts ssh.ConnectOptions) error {
	server, err := a.Config.FindServer(serverName)
	if err != nil {
		return err
	}

	// Apply defaults
	resolved := a.resolveServer(server)

	// Apply overrides
	if opts.User != "" {
		resolved.User = opts.User
	}
	if opts.Port != 0 {
		resolved.Port = opts.Port
	}
	if opts.Key != "" {
		resolved.Key = opts.Key
	}

	return a.SSHClient.Connect(resolved)
}

// resolveServer applies defaults to a server config
func (a *App) resolveServer(s *config.Server) *config.Server {
	resolved := *s

	if resolved.User == "" && a.Config.Defaults.User != "" {
		resolved.User = a.Config.Defaults.User
	}
	if resolved.Port == 0 {
		if a.Config.Defaults.Port != 0 {
			resolved.Port = a.Config.Defaults.Port
		} else {
			resolved.Port = 22
		}
	}
	if resolved.Key == "" && a.Config.Defaults.Key != "" {
		resolved.Key = a.Config.Defaults.Key
	}

	return &resolved
}

// Save persists the config to disk
func (a *App) Save() error {
	return a.Config.Save()
}
