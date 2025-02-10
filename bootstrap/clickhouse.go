package bootstrap

import (
	"context"
	"fmt"
	"shalabing-gin/global"
	"time"

	// "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

func InitializeClickHouse() clickhouse.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", global.App.Config.Clickhouse.Host, global.App.Config.Clickhouse.Port)},
		Auth: clickhouse.Auth{
			Database: global.App.Config.Clickhouse.Database,
			Username: global.App.Config.Clickhouse.Username,
			Password: global.App.Config.Clickhouse.Password,
		},
		Debug: true,
	})
	if err != nil {
		global.App.Log.Error("Can not connect ClickHouse: ", zap.Any("err", err))
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		global.App.Log.Error("ClickHouse connect ping failed:", zap.Any("err", err))
	}
	return conn
}
