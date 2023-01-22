package redix

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

func GetProto(ctx context.Context, rdb *redis.ClusterClient, message proto.Message, key string) error {
	protoBytes, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	err = proto.Unmarshal(protoBytes, message)
	if err != nil {
		return fmt.Errorf("unmarshal [%s] proto failed, %s", message.ProtoReflect().Descriptor().Name(), err.Error())
	}

	return nil
}

func SetProto(ctx context.Context, rdb *redis.ClusterClient, message proto.Message, key string, expiration time.Duration) error {
	protoBytes, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal [%s] proto failed, %s", message.ProtoReflect().Descriptor().Name(), err.Error())
	}

	err = rdb.Set(ctx,
		key,
		protoBytes,
		expiration,
	).Err()

	return err
}

func DelMatchedKeysInMaster(ctx context.Context, rdb *redis.ClusterClient, scanNum int64, match string) error {
	err := rdb.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		var keys []string
		cursor := uint64(0)
		var err error

		wg := &sync.WaitGroup{}

		for {
			keys, cursor, err = client.Scan(ctx, cursor, match, scanNum).Result()
			if err != nil {
				return err
			}

			for _, key := range keys {
				wg.Add(1)
				go func(k string) {
					rdb.Del(ctx, k)
					wg.Done()
				}(key)
			}

			if cursor == 0 {
				break
			}
		}

		wg.Wait()

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func Transaction(ctx context.Context, rdb *redis.ClusterClient, key string, txf func(tx *redis.Tx) error) error {
	for i := 0; i < 10; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			return nil
		}
		if err == redis.TxFailedErr {
			time.Sleep(time.Second)
			continue
		}
		return err
	}

	return errors.New("transaction reached maximum number of retries")
}
