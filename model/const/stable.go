package stable

// 状态：0.初始状态 1.草稿状态 2.生效中 3.已删除
type RecordStatus int

const (
	// StatusInitial 初始状态
	StatusInitial = iota
	// StatusDraft 草稿状态
	StatusDraft
	// StatusActive 生效中
	StatusActive
	// StatusDeleted 已删除
	StatusDeleted
)

var RecordStatusMap = map[int]string{
	StatusInitial: "初始状态",
	StatusDraft:   "草稿",
	StatusActive:  "已完成",
	StatusDeleted: "已删除",
}
