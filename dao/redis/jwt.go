package redis

import (
	"bluebell_backend/settings"
	"time"
)

// Redis服务器添加token设置TTL
func AddToken(userid string, authToken string) error {
	_, err := client.Set(userid, authToken, time.Duration(settings.Conf.TokenExpireDuration)*time.Second).Result()
	return err
}

// 获取Token
func GetToken(userid string) (string, error) {
	authToken, err := client.Get(userid).Result()
	return authToken, err
}
