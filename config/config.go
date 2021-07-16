// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"io/ioutil"

	pconfig "github.com/prometheus/common/config"
	"gopkg.in/yaml.v2"
)

// Metric contains values that define a metric
type File struct {
	Path    string            `yaml:"path"`
	Metrics []Metric          `yaml:"metrics"`
	Headers map[string]string `yaml:"headers,omitempty"`
}
type Metric struct {
	Name            string
	Path            string
	Timestamp       string
	TimestampFormat string
	Timezone        string
	Labels          map[string]string
	Type            MetricType
	Help            string
	Values          map[string]string
}

type MetricType string

const (
	ValueScrape  MetricType = "value" // default
	ObjectScrape MetricType = "object"
)

// Config contains metrics and headers defining a configuration
type Config struct {
	Files            []File                   `yaml:"files"`
	HTTPClientConfig pconfig.HTTPClientConfig `yaml:"http_client_config,omitempty"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	for _, file := range config.Files {

		// Complete Defaults
		for i := 0; i < len(file.Metrics); i++ {
			if file.Metrics[i].Type == "" {
				file.Metrics[i].Type = ValueScrape
			}
			if file.Metrics[i].Help == "" {
				file.Metrics[i].Help = file.Metrics[i].Name
			}
			if file.Metrics[i].Timezone == "" {
				file.Metrics[i].Timezone = "UTC"
			}
			if file.Metrics[i].TimestampFormat == "" {
				file.Metrics[i].TimestampFormat = "02.01.2006 15:04:05"
			}
		}

	}

	return config, nil
}
