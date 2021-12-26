package pulseaudio

type SubscriptionEventType uint32

const (
	SubscriptionEventSink SubscriptionEventType = iota
	SubscriptionEventSource
	SubscriptionEventSourceInput
	SubscriptionEventSourceOutput
	SubscriptionEventModule
	SubscriptionEventClient
	SubscriptionEventSimpleCache
	SubscriptionEventServer
	SubscriptionEventAutoload
	SubscriptionEventCard

	SubscriptionEventFacilityMask = 0x000f
	SubscriptionEventTypeNew      = 0x0000
	SubscriptionEventTypeChange   = 0x0010
	SubscriptionEventTypeRemove   = 0x0020
	SubscriptionEventTypeMask     = 0x0030
)
