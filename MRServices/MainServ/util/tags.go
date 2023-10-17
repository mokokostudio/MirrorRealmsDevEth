package util

const (
	TagsKeyToken         = "token"
	TagsKeyIPRegion      = "ip_region"
	TagsKeyUserTagging   = "user_tagging"
	TagsKeyAccount       = "account"
	TagsKeyRegion        = "region"
	TagsKeyUserID        = "user_id"
	TagsKeyClientVersion = "client_version"
)

type UserTagging struct {
	TaggingName string `json:"tagging_name"`
	Tagging     string `json:"tagging"`
}
