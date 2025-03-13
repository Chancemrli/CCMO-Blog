package redis

/*
	Redis Key
*/

const (
	KeyPostInfoHashPrefix = "ccmo:post:"
	KeyPostTimeZSet       = "ccmo:post:time"
	KeyPostScoreZSet      = "ccmo:post:score"
	//KeyPostVotedUpSetPrefix   = "ccmo:post:voted:down:"
	//KeyPostVotedDownSetPrefix = "ccmo:post:voted:up:"
	KeyPostVotedZSetPrefix = "ccmo:post:voted:"

	KeyCommunityPostSetPrefix = "ccmo:community:"

	KeyHotArticlesMap   = "cache:post:hot"
	KeyHotArticlePrefix = "cache:post:hot:"

	MinExpireTime = 300 // seconds
	MaxExpireTime = 480 // seconds
)
