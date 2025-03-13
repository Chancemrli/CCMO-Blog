package redis

func SetEmptyArticle(postID string) (bool, error) {
	key := KeyHotArticlePrefix + postID
	ttl := RandTTL()
	res, err := client.HSet(key, "post_id", 0).Result()
	client.Expire(key, ttl)
	return res, err
}
