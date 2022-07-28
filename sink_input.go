package pulseaudio

import "io"

// SinkInput contains information about a sink inputs in pulseaudio
// https://github.com/pulseaudio/pulseaudio/blob/7f4d7fcf5f6407913e50604c6195d0d5356195b1/src/pulsecore/protocol-native.c#L3397
// https://freedesktop.org/software/pulseaudio/doxygen/structpa__sink__input__info.html
type SinkInput struct {
	Index          uint32
	Name           string
	ModuleIndex    uint32
	ClientIndex    uint32
	SinkIndex      uint32
	SampleSpec     sampleSpec
	ChannelMap     channelMap
	Cvolume        cvolume
	BufferLatency  uint64
	SinkLatency    uint64
	ResampleMethod string
	Driver         string
	Muted          bool
	PropList       map[string]string
	Corked         bool
	HasVolume      bool
	VolumeWritable bool
	Format         formatInfo
}

// ReadFrom deserializes a sink input packet from pulseaudio
func (s *SinkInput) ReadFrom(r io.Reader) (int64, error) {
	err := bread(r,
		uint32Tag, &s.Index,
		stringTag, &s.Name,
		uint32Tag, &s.ModuleIndex,
		uint32Tag, &s.ClientIndex,
		uint32Tag, &s.SinkIndex,
		&s.SampleSpec,
		&s.ChannelMap,
		&s.Cvolume,
		usecTag, &s.BufferLatency,
		usecTag, &s.SinkLatency,
		stringTag, &s.ResampleMethod,
		stringTag, &s.Driver,
		&s.Muted,
		&s.PropList,
		&s.Corked,
		&s.HasVolume,
		&s.VolumeWritable,
		&s.Format,
	)
	return 0, err
}

// SinkInputs queries PulseAudio for a list of sink inputs and returns an array
func (c *Client) SinkInputs() ([]SinkInput, error) {
	b, err := c.request(commandGetSinkInputInfoList)
	if err != nil {
		return nil, err
	}
	var sinkInputs []SinkInput
	for b.Len() > 0 {
		var sinkInput SinkInput
		err = bread(b, &sinkInput)
		if err != nil {
			return nil, err
		}
		sinkInputs = append(sinkInputs, sinkInput)
	}
	return sinkInputs, nil
}

func (c *Client) GetSinkInputInfo(index uint32) (*SinkInput, error) {
	b, err := c.request(commandGetSinkInputInfo, uint32Tag, index)
	if err != nil {
		return nil, err
	}
	var sinkInput SinkInput
	err = bread(b, &sinkInput)
	if err != nil {
		return nil, err
	}
	return &sinkInput, nil
}

func (c *Client) SetSinkInputVolume(index uint32, volume float32) error {
	return c.setSinkInputVolume(index, cvolume{uint32(volume * 0xffff)})
}

func (c *Client) setSinkInputVolume(index uint32, cvolume cvolume) error {
	// https://github.com/pulseaudio/pulseaudio/blob/7f4d7fcf5f6407913e50604c6195d0d5356195b1/src/pulsecore/protocol-native.c#L3760
	_, err := c.request(commandSetSinkInputVolume, uint32Tag, index, cvolume)
	return err
}

func (c *Client) SetSinkInputMute(index uint32, mute bool) error {
	muteCmd := '0'
	if mute {
		muteCmd = '1'
	}

	_, err := c.request(commandSetSinkInputMute, uint32Tag, index, uint8(muteCmd))
	return err
}
