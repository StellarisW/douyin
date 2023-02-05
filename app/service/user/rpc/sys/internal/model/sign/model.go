package sign

import (
	"context"
	"douyin/app/common/errx"
	"douyin/app/common/log"
	"douyin/app/service/user/internal/user"
	"douyin/app/service/user/rpc/sys/internal/model/dao/entity"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/yitter/idgenerator-go/idgen"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

const clientId = "douyin"

type (
	Model interface {
		Register(ctx context.Context, username, password string) (int64, string, errx.Error)
		Login(ctx context.Context, username, password string) (int64, string, errx.Error)
	}
	DefaultModel struct {
		idGenerator *idgen.DefaultIdGenerator
		db          *gorm.DB
		rdb         *redis.ClusterClient

		AuthString string
	}
)

func NewModel(v *viper.Viper, idGenerator *idgen.DefaultIdGenerator, db *gorm.DB, rdb *redis.ClusterClient) *DefaultModel {
	clientSecret := v.GetString("Client." + clientId + ".Secret")
	if clientSecret == "" {
		log.Logger.Fatal("get client secret failed")
	}

	encodeAuthString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)))
	basicAuthString := "Basic " + encodeAuthString

	return &DefaultModel{
		idGenerator: idGenerator,
		db:          db,
		rdb:         rdb,

		AuthString: basicAuthString,
	}
}

func (m *DefaultModel) Register(ctx context.Context, username, password string) (int64, string, errx.Error) {
	err := m.db.WithContext(ctx).
		Table(entity.TableNameUserSubject).
		Select("`id`").
		Where("`username` = ?", username).
		Take(&entity.UserSubject{}).Error
	switch err {
	case nil:
		// 用户已经存在的情况

		return 0, "", errUsernameExists

	case gorm.ErrRecordNotFound:
		// 用户不存在的情况(可以创建用户)

		userId := m.idGenerator.NewLong()

		err = m.db.WithContext(ctx).
			Create(&entity.UserSubject{
				ID:       userId,
				Username: username,
				Password: encryptPassword(password),
			}).Error
		if err != nil {
			log.Logger.Error(errx.MysqlInsert, zap.Error(err))
			return 0, "", errMysqlInsert
		}

		token, erx := m.getToken(ctx, userId)
		if erx != nil {
			return 0, "", erx
		}

		// 添加注册缓存
		err = m.rdb.SAdd(ctx,
			user.RdbKeyRegisterSet,
			username).Err()
		if err != nil {
			log.Logger.Error(errx.RedisSet, zap.Error(err))
			return 0, "", errRedisSet
		}

		return userId, token, nil
	default:
		// 数据库查询失败的情况

		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return 0, "", errMysqlGet
	}
}

func (m *DefaultModel) Login(ctx context.Context, username, password string) (int64, string, errx.Error) {
	var userId int64

	err := m.db.WithContext(ctx).
		Table(entity.TableNameUserSubject).
		Select("`id`").
		Where("`username` = ? AND `password` = ?", username, encryptPassword(password)).
		Take(&userId).
		Error
	switch err {
	case nil:
		// 查看账户登录次数
		cnt, err := m.rdb.Get(ctx,
			user.RdbKeyLoginFrozenLoginCnt+username).Int64()
		if err != nil {
			if err != redis.Nil {
				log.Logger.Error(errx.RedisGet, zap.Error(err))
				return 0, "", errRedisGet
			}
		} else {
			if cnt > 10 {
				// 短时间内登录大于10次则冻结
				lastFrozenTime, err := m.rdb.Get(ctx,
					user.RdbKeyLoginFrozenTimeLast+username).Result()
				if err != nil {
					if err != redis.Nil {
						log.Logger.Error(errx.RedisGet, zap.Error(err))
						return 0, "", errRedisGet
					}
				}

				var lastDuration time.Duration
				var currentDuration time.Duration
				switch lastFrozenTime {
				case "":
					lastDuration, _ = time.ParseDuration("1s")
					currentDuration, _ = time.ParseDuration("1m")
				case "1s":
					lastDuration, _ = time.ParseDuration("1m")
					currentDuration, _ = time.ParseDuration("15m")
				case "1m":
					lastDuration, _ = time.ParseDuration("15m")
					currentDuration, _ = time.ParseDuration("1h")
				case "15m":
					lastDuration, _ = time.ParseDuration("1h")
					currentDuration, _ = time.ParseDuration("24h")
				case "1h":
					lastDuration, _ = time.ParseDuration("24h")
					currentDuration, _ = time.ParseDuration("168h")
				case "1d":
					lastDuration, _ = time.ParseDuration("168h")
					currentDuration, _ = time.ParseDuration("720h")
				case "7d":
					lastDuration = -1
					currentDuration = -1
				}

				err = m.rdb.Set(ctx,
					user.RdbKeyLoginFrozenTimeLast+username,
					1,
					lastDuration).Err()
				if err != nil {
					log.Logger.Error(errx.RedisGet, zap.Error(err))
					return 0, "", errRedisGet
				}

				err = m.rdb.Set(ctx,
					user.RdbKeyLoginFrozenTime+username,
					1,
					currentDuration).Err()
				if err != nil {
					log.Logger.Error(errx.RedisSet, zap.Error(err))
					return 0, "", errRedisSet
				}

				err = m.rdb.Del(ctx,
					user.RdbKeyLoginFrozenLoginCnt+username).Err()
				if err != nil {
					log.Logger.Error(errx.RedisDel, zap.Error(err))
					return 0, "", errRedisDel
				}
			}
		}

		// 登录次数增加一次
		cnt, err = m.rdb.Incr(ctx,
			user.RdbKeyLoginFrozenLoginCnt+username).Result()
		if err != nil {
			log.Logger.Error(errx.RedisIncr, zap.Error(err))
			return 0, "", errRedisIncr
		} else {
			if cnt == 1 {
				m.rdb.Expire(ctx,
					user.RdbKeyLoginFrozenLoginCnt+username,
					time.Minute*15)
			}
		}

		token, erx := m.getToken(ctx, userId)
		if erx != nil {
			return 0, "", erx
		}

		return userId, token, nil
	case gorm.ErrRecordNotFound:

		return 0, "", errWrongUsernameOrPassword

	default:

		log.Logger.Error(errx.MysqlGet, zap.Error(err))
		return 0, "", errMysqlGet
	}
}

func (m *DefaultModel) getToken(ctx context.Context, userId int64) (string, errx.Error) {
	authRes, err := req.NewRequest().
		SetHeader("Authorization", m.AuthString).
		SetQueryParam("obj", strconv.FormatInt(userId, 10)).
		Get("http://douyin-auth-api:11120/douyin/token/auth")
	if err != nil {
		log.Logger.Error(errx.RequestHttpSend, zap.Error(err))
		return "", errRequestHttpSend
	}
	if authRes.StatusCode != http.StatusOK {
		log.Logger.Error(errx.RequestHttpStatusCode, zap.Error(err), zap.String("status", authRes.Status))
		return "", errRequestHttpStatusCode
	}

	authResJson := gjson.Parse(authRes.String())

	token := authResJson.Get("data.token.token_value").String()

	return token, nil
}

func encryptPassword(password string) string {
	d := sha3.Sum224([]byte(password))
	return hex.EncodeToString(d[:])
}
