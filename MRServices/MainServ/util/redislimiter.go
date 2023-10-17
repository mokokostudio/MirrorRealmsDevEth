package util

// type RedisLimiter struct {
// 	bot    gredis.RedisClient
// 	logger    *zap.Logger
// 	duration  int64 // duration is in milliseconds
// 	cntPerDur int64 // count per duration
// }
//
// func NewRedisLimiter(bot gredis.RedisClient, logger *zap.Logger, duration, cntPerDur int64) (*RedisLimiter, error) {
// 	return &RedisLimiter{bot: bot, logger: logger, duration: duration * 1e3, cntPerDur: cntPerDur}, nil
// }
//
// func (rl *RedisLimiter) Can(ctx context.Context, key string) bool {
// 	datas, err := ToInt64s(rl.bot.ZRange(ctx, key, 0, -1))
// 	if err != nil {
// 		rl.logger.Error("can not read limiter data from gredis! ", zap.String("key", key),
// 			zap.Int64("duration", rl.duration), zap.Int64("cntPerDur", rl.cntPerDur), zap.Error(err))
// 		return true
// 	}
// 	nowms := NowInMilliSeconds()
// 	tm := nowms - rl.duration
//
// 	var i int
// 	for _, v := range datas {
// 		if v < tm {
// 			i++
// 			continue
// 		}
// 		break
// 	}
//
// 	datas = datas[i:]
//
// 	if len(datas) >= int(rl.cntPerDur) {
// 		return false
// 	}
// 	_, err = rl.bot.ZAdd(ctx, key, float64(nowms), strconv.Itoa(int(nowms)))
// 	if err != nil {
// 		rl.logger.Error("can not write limiter data to gredis! ", zap.String("key", key),
// 			zap.Int64("duration", rl.duration), zap.Int64("cntPerDur", rl.cntPerDur), zap.Error(err))
// 	}
// 	_, err = rl.bot.Expire(ctx, key, int(rl.duration/1e3))
// 	if err != nil {
// 		rl.logger.Error("can not write limiter data to gredis! ", zap.String("key", key),
// 			zap.Int64("duration", rl.duration), zap.Int64("cntPerDur", rl.cntPerDur), zap.Error(err))
// 	}
// 	return true
// }
