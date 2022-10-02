package consts

const (
	TaskTelegramPrefix = "telegram:"
	TaskTelegramUpdate = TaskTelegramPrefix + "update"
)

const (
	TaskFeedPrefix    = "feed:"
	TaskFeedScheduler = TaskFeedPrefix + "scheduler"
	TaskFeedImporter  = TaskFeedPrefix + "importer"
	TaskFeedSave      = TaskFeedPrefix + "save"
	TaskFeedAdd       = TaskFeedPrefix + "add"
	TaskFeedView      = TaskFeedPrefix + "view"
)

const (
	TaskFeedItemPrefix     = "feedItem:"
	TaskFeedItemSave       = TaskFeedItemPrefix + "save"
	TaskFeedItemView       = TaskFeedItemPrefix + "view"
	TaskFeedItemGroup      = TaskFeedItemPrefix + "group"
	TaskFeedItemAggregated = TaskFeedItemPrefix + "aggregated"
	TaskFeedItemBroadcast  = TaskFeedItemPrefix + "broadcast"
)

const (
	TaskRoomPrefix = "room:"
	TaskRoomSave   = TaskRoomPrefix + "save"
	TaskRoomAdd    = TaskRoomPrefix + "add"
	TaskRoomLeft   = TaskRoomPrefix + "left"
	TaskRoomView   = TaskRoomPrefix + "view"
)
