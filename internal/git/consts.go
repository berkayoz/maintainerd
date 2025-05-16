package git

import "fmt"

var workDir = "/tmp"

var gitConfigPath = fmt.Sprintf("%s/.gitconfig", workDir)

var envOptions = []string{fmt.Sprintf("GIT_CONFIG_GLOBAL=%s", gitConfigPath), fmt.Sprintf("GIT_CONFIG_SYSTEM=%s", gitConfigPath)}
