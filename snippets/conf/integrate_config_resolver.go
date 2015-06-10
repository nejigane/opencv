package conf

import (
	"encoding/json"
	"io/ioutil"
	"pfi/InStoreAutomation/kanohi-scouter-conf"
)

// IntegrateConfig is parameters of Integrate snippets.
type IntegrateConfig struct {
	IntegratorConfig         string
	InstanceManagerConfig    string
	VisualizerConfig         string
	FrameInputKeys           []string
	DetectionResultInputKeys []string
	FloorID                  int
	PlayerFlag               bool
	JpegQuality              int
}

// GetIntegrateConfig crates configuration data reading external file.
func GetIntegrateConfig(filePath string) (IntegrateConfig, error) {
	conf := IntegrateConfig{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return conf, err
	}

	var integrateConfig ksconf.IntegrateSnippet
	err = json.Unmarshal(file, &integrateConfig)
	if err != nil {
		return conf, err
	}

	// get scouter::Integrate::Config
	integratorConfByte, err := json.Marshal(integrateConfig.Integrator)
	if err != nil {
		return conf, err
	}
	integratorConf := string(integratorConfByte)

	instanceManagerByte, err := json.Marshal(integrateConfig.InstanceManager)
	if err != nil {
		return conf, err
	}
	instanceManagerConf := string(instanceManagerByte)

	visualizerByte, err := json.Marshal(integrateConfig.Visualizer)
	if err != nil {
		return conf, err
	}
	visualizerConf := string(visualizerByte)

	return IntegrateConfig{
		IntegratorConfig:         integratorConf,
		InstanceManagerConfig:    instanceManagerConf,
		VisualizerConfig:         visualizerConf,
		FrameInputKeys:           integrateConfig.FrameInputKeys,
		DetectionResultInputKeys: integrateConfig.DetectionResultInputKeys,
		FloorID:                  integrateConfig.FloorID,
		PlayerFlag:               integrateConfig.Player != nil,
		JpegQuality:              50,
	}, nil
}
