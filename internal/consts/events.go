package consts

const (
	EventErrorPrefix    = "error"
	EventErrorForbidden = EventErrorPrefix + ":forbidden"
	EventErrorNotFound  = EventErrorPrefix + ":notfound"
	EventErrorViewList  = EventErrorPrefix + ":viewlist"
	EventErrorArgs      = EventErrorPrefix + ":args"
)

const (
	EventFeedPrefix    = "feed"
	EventFeedSave      = EventFeedPrefix + ":save"
	EventFeedSaveError = EventFeedSave + ":error"
	EventFeedSaveAfter = EventFeedSave + ":after"
	EventFeedView      = EventFeedPrefix + ":view"
	EventFeedViewOne   = EventFeedView + ":one"
	EventFeedViewList  = EventFeedView + ":list"
)

const (
	EventFeedItemPrefix   = "feedItem"
	EventFeedItemView     = EventFeedItemPrefix + ":view"
	EventFeedItemViewList = EventFeedItemView + ":list"
)

const (
	EventRoomPrefix    = "room"
	EventRoomSave      = EventRoomPrefix + ":save"
	EventRoomSaveError = EventRoomSave + ":error"
	EventRoomSaveAfter = EventRoomSave + ":after"
	EventRoomView      = EventRoomPrefix + ":view"
	EventRoomViewOne   = EventRoomView + ":one"
	EventRoomViewList  = EventRoomView + ":list"
)
