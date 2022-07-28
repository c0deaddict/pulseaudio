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

	EventFacilityMask Event = 0x000f
	EventTypeNew      Event = 0x0000
	EventTypeChange   Event = 0x0010
	EventTypeRemove   Event = 0x0020
	EventTypeMask     Event = 0x0030
)

type SubscriptionEvent struct {
	Event Event
	Index uint32
}

// Updates returns a channel with PulseAudio updates.
func (c *Client) Updates() (updates <-chan SubscriptionEvent, err error) {
	const subscriptionMaskAll = 0x02ff
	_, err = c.request(commandSubscribe, uint32Tag, uint32(subscriptionMaskAll))
	if err != nil {
		return nil, err
	}
	return c.updates, nil
}
