package topic

import (
	"admin_base_server/api/v1/topic"
)

type RouterGroup struct {
	TopicRouter
}

var (
	topicAPI      *topic.TopicAPI
	topicGroupAPI *topic.TopicGroupAPI
)
