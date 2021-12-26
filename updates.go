package pulseaudio

type Event uint32

const (
	EventSink Event = iota
	EventSource
	EventSinkInput
	EventSourceOutput
	EventModule
	EventClient
	EventSimpleCache
	EventServer
	EventAutoload
	EventCard

	EventFacilityMask = 0x000f
	EventTypeNew      = 0x0000
	EventTypeChange   = 0x0010
	EventTypeRemove   = 0x0020
	EventTypeMask     = 0x0030
)

// Updates returns a channel with PulseAudio updates.
func (c *Client) Updates() (updates <-chan Event, err error) {
	const subscriptionMaskAll = 0x02ff
	_, err = c.request(commandSubscribe, uint32Tag, uint32(subscriptionMaskAll))
	if err != nil {
		return nil, err
	}
	return c.updates, nil
}
