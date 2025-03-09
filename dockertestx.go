package dockertestx

import "github.com/ory/dockertest/v3"

// RunOption is a function that modifies a dockertest.RunOptions.
type RunOption func(*dockertest.RunOptions)
