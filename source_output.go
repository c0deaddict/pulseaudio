package pulseaudio

import "io"

// SinkInput contains information about a sink inputs in pulseaudio
// https://freedesktop.org/software/pulseaudio/doxygen/structpa__source__output__info.html
// https://github.com/pulseaudio/pulseaudio/blob/7f4d7fcf5f6407913e50604c6195d0d5356195b1/src/pulsecore/protocol-native.c#L3440
type SourceOutput struct {
	Index          uint32
	Name           string
	ModuleIndex    uint32
	ClientIndex    uint32
	SourceIndex    uint32
	SampleSpec     sampleSpec
	ChannelMap     channelMap
	BufferLatency  uint64
	SourceLatency  uint64
	ResampleMethod string
	Driver         string
	PropList       map[string]string
	Corked         bool
	Cvolume        cvolume
	Muted          bool
	HasVolume      bool
	VolumeWritable bool
	Format         formatInfo
}

// ReadFrom deserializes a sink input packet from pulseaudio
func (s *SourceOutput) ReadFrom(r io.Reader) (int64, error) {
	err := bread(r,
		uint32Tag, &s.Index,
		stringTag, &s.Name,
		uint32Tag, &s.ModuleIndex,
		uint32Tag, &s.ClientIndex,
		uint32Tag, &s.SourceIndex,
		&s.SampleSpec,
		&s.ChannelMap,
		usecTag, &s.BufferLatency,
		usecTag, &s.SourceLatency,
		stringTag, &s.ResampleMethod,
		stringTag, &s.Driver,
		&s.PropList,
		&s.Corked,
		&s.Cvolume,
		&s.Muted,
		&s.HasVolume,
		&s.VolumeWritable,
		&s.Format,
	)
	return 0, err
}

// SinkInputs queries PulseAudio for a list of sink inputs and returns an array
func (c *Client) SourceOutputs() ([]SourceOutput, error) {
	b, err := c.request(commandGetSourceOutputInfoList)
	if err != nil {
		return nil, err
	}
	var sourceOutputs []SourceOutput
	for b.Len() > 0 {
		var sourceOutput SourceOutput
		err = bread(b, &sourceOutput)
		if err != nil {
			return nil, err
		}
		sourceOutputs = append(sourceOutputs, sourceOutput)
	}
	return sourceOutputs, nil
}

func (c *Client) GetSourceOutputInfo(index uint32) (*SourceOutput, error) {
	b, err := c.request(commandGetSourceOutputInfo, uint32Tag, index)
	if err != nil {
		return nil, err
	}
	var sourceOutput SourceOutput
	err = bread(b, &sourceOutput)
	if err != nil {
		return nil, err
	}
	return &sourceOutput, nil
}

func (c *Client) SetSourceOutputVolume(index uint32, volume float32) error {
	return c.setSourceOutputVolume(index, cvolume{uint32(volume * 0xffff)})
}

func (c *Client) setSourceOutputVolume(index uint32, cvolume cvolume) error {
	// https://github.com/pulseaudio/pulseaudio/blob/7f4d7fcf5f6407913e50604c6195d0d5356195b1/src/pulsecore/protocol-native.c#L3760
	_, err := c.request(commandSetSourceOutputVolume, uint32Tag, index, cvolume)
	return err
}

func (c *Client) SetSourceOutputMute(index uint32, mute bool) error {
	muteCmd := '0'
	if mute {
		muteCmd = '1'
	}

	_, err := c.request(commandSetSourceOutputMute, uint32Tag, index, uint8(muteCmd))
	return err
}
