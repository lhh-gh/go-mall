package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github/lhh-gh/go-mall/comon/enum"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/logic/do"
)

// 关于使用 HSET 直接存储结构体的讨论 https://github.com/redis/go-redis/discussions/2454

type DummyDemoOrder struct {
	OrderNo string `redis:"orderNo"`
	UserId  int64  `redis:"userId"`
}

// SetDemoOrderStruct 使用HSET的存储结构体数据
func SetDemoOrderStruct(ctx context.Context, demoOrder *do.DemoOrder) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	// 构造 Redis Key，格式为预定义模板 + 订单编号

	data := struct { // 创建一个匿名结构体用于存储到 Redis
		OrderNo string `redis:"orderNo"`
		UserId  int64  `redis:"userId"`
	}{
		UserId:  demoOrder.UserId,
		OrderNo: demoOrder.OrderNo,
	}

	_, err := Redis().HSet(ctx, redisKey, data).Result() // 使用 Redis HSET 命令存储结构体数据
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err) // 如果出错，记录错误日志
		return err
	}

	return nil
}

// GetDemoOrderStruct 使用HGETALL 和 Scan 读取结构体数据
func GetDemoOrderStruct(ctx context.Context, orderNo string) (*DummyDemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	// 构造 Redis Key，格式为预定义模板 + 订单编号

	data := new(DummyDemoOrder) // 初始化一个 DummyDemoOrder 结构体指针
	err := Redis().HGetAll(ctx, redisKey).Scan(&data)
	// 从 Redis 中获取所有字段并映射到结构体

	Redis().Get(ctx, redisKey).String() // 此行未使用结果，可能为冗余代码

	if err != nil {
		logger.New(ctx).Error("redis error", "err", err) // 如果出错，记录错误日志
		return nil, err
	}

	logger.New(ctx).Info("scan data from redis", "data", &data) // 输出读取到的数据信息
	return data, nil
}

func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	jsonDataBytes, _ := json.Marshal(demoOrder)
	// 将 DemoOrder 对象序列化为 JSON 字节数组

	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	// 构造 Redis Key，格式为预定义模板 + 订单编号

	_, err := Redis().Set(ctx, redisKey, jsonDataBytes, 0).Result()
	// 使用 Redis SET 命令将 JSON 数据存入 Redis

	if err != nil {
		logger.New(ctx).Error("redis error", "err", err) // 如果出错，记录错误日志
		return err
	}

	return nil
}

/**
  日志链路追踪
*/

func GetDemoOrder(ctx context.Context, orderNo string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	// 构造 Redis Key，格式为预定义模板 + 订单编号

	jsonBytes, err := Redis().Get(ctx, redisKey).Bytes()
	// 从 Redis 获取对应的字节数据

	if err != nil {
		logger.New(ctx).Error("redis error", "err", err) // 如果出错，记录错误日志
		return nil, err
	}

	data := new(do.DemoOrder)        // 初始化一个新的 DemoOrder 对象
	json.Unmarshal(jsonBytes, &data) // 反序列化 JSON 字节数据到对象中
	return data, nil
}
