package cache

type Channel string

const (
	PChKeyEventsExpire Channel = "__keyevent@*__:expired"
	PChKeyEventsHSet   Channel = "__keyevent@*__:hset"
	PChKeyEventsHDel   Channel = "__keyevent@*__:hdel"
	PChKeyEventsDel    Channel = "__keyevent@*__:del"
)
