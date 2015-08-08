package detector

import (
	"io/ioutil"
	"pfi/sensorbee/scouter/bridge"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

// FrameProcessorParamState is a shared state of frame processor parameter for
// scouter-core.
type FrameProcessorParamState struct {
	fp bridge.FrameProcessor
}

func createFrameProcessorParamState(ctx *core.Context, params data.Map) (core.SharedState,
	error) {
	p, err := params.Get("file")
	if err != nil {
		return nil, err
	}
	path, err := data.AsString(p)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fpConfig := string(b)
	s := &FrameProcessorParamState{}
	s.fp = bridge.NewFrameProcessor(fpConfig)

	return s, nil
}

// CreateNewState creates a state of frame processor parameters. The parameter
// is collected on JSON file, see `scouter::FrameProcessor::Config`, which is
// composition of camera parameters and RIO information.
// This state is updatable.
//
// Usage of WITH parameters:
//  file: The file path. Returns an error when cannot read the file.
func (s *FrameProcessorParamState) CreateNewState() func(*core.Context, data.Map) (
	core.SharedState, error) {
	return createFrameProcessorParamState
}

func (s *FrameProcessorParamState) TypeName() string {
	return "scouter_frame_processor_param"
}

func (s *FrameProcessorParamState) Terminate(ctx *core.Context) error {
	s.fp.Delete()
	return nil
}

// Update the state to reload the JSON file without lock.
//
// Usage of WITH parameters:
//  file: The file path. Returns an error when cannot read the file.
func (s *FrameProcessorParamState) Update(params data.Map) error {
	p, err := params.Get("file")
	if err != nil {
		return err
	}
	path, err := data.AsString(p)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fpConfig := string(b)
	s.fp.UpdateConfig(fpConfig)

	return nil
}