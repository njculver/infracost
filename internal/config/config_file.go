package config

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

const minConfigFileVersion = "0.1"
const maxConfigFileVersion = "0.1"

type ConfigFileSpec struct { // nolint:revive
	Version  string     `yaml:"version"`
	Projects []*Project `yaml:"projects" ignored:"true"`
	Config   []*Config  `yaml:"config" ignored:"true"`
}

func LoadConfigFile(path string) (ConfigFileSpec, error) {
	cfgFile := ConfigFileSpec{}

	if !FileExists(path) {
		return cfgFile, errors.Errorf("Config file does not exist at %s", path)
	}

	rawCfgFile, err := ioutil.ReadFile(path)
	if err != nil {
		return cfgFile, err
	}

	err = yaml.Unmarshal(rawCfgFile, &cfgFile)
	if err != nil {
		return cfgFile, errors.New("Error parsing config YAML: " + strings.TrimPrefix(err.Error(), "yaml: "))
	}

	if !checkVersion(cfgFile.Version) {
		return cfgFile, errors.Errorf("Invalid config file version. Supported versions are %s ≤ x ≤ %s", minConfigFileVersion, maxConfigFileVersion)
	}

	return cfgFile, nil
}

func checkVersion(v string) bool {
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return semver.Compare(v, "v"+minConfigFileVersion) >= 0 && semver.Compare(v, "v"+maxConfigFileVersion) <= 0
}
