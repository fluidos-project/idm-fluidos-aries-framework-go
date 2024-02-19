package files

import _ "embed" // required for tests only

// Sample testdata files to be used for tests only.
// nolint:gochecknoglobals
var (
	//go:embed rawVC.json
	ExampleRawVC []byte
	//go:embed sampleFrame.json
	SampleFramePsms []byte
	//go:embed sampleFrameFrame.json
	SampleFramePsmsFrame []byte
	//go:embed sampleFrameFrameOther.json
	SampleFramePsmsFrameOther []byte
)
