package global

import (
	"shalabing-gin/config"

	// "github.com/go-redis/redis"

	"shalabing-gin/go-storage/storage"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	Log         *zap.Logger
	DB          *gorm.DB
	Redis       *redis.ClusterClient //*redis.Client
	Cron        *cron.Cron
	Rabbitmq    *Queue
	Kafka       *Queue
	MongoDB     *mongo.Client //*mongo.Database
	ES          *elasticsearch.Client
	Clickhouse  clickhouse.Conn
}

var App = new(Application)

func (app *Application) Disk(disk ...string) storage.Storage {
	diskName := app.Config.Storage.Default
	if len(disk) > 0 {
		diskName = storage.DiskName(disk[0])
	}
	s, err := storage.Disk(diskName)
	if err != nil {
		panic(err)
	}
	return s
}
