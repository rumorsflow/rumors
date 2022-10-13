package consts

const TaskRumors = "rumors"

const (
	TaskFeedPrefix   = "feed:"
	TaskFeedImporter = TaskFeedPrefix + "importer"
)

const (
	TaskFeedItemPrefix    = "feedItem:"
	TaskFeedItemGroup     = TaskFeedItemPrefix + "group"
	TaskFeedItemAggregate = TaskFeedItemPrefix + "aggregate"
)

const (
	TaskRoomPrefix    = "room:"
	TaskRoomBroadcast = TaskRoomPrefix + "broadcast"
	TaskRoomStart     = TaskRoomPrefix + "start"
	TaskRoomUpdated   = TaskRoomPrefix + "updated"
)

const (
	TaskTelegramPrefix = "telegram:"
	TaskTelegramUpdate = TaskTelegramPrefix + "update"
)
