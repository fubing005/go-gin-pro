package bootstrap

import (
	"shalabing-gin/global"

	"shalabing-gin/go-storage/kodo"
	"shalabing-gin/go-storage/local"
	"shalabing-gin/go-storage/oss"
)

func InitializeStorage() {
	_, _ = local.Init(global.App.Config.Storage.Disks.Local)
	_, _ = kodo.Init(global.App.Config.Storage.Disks.QiNiu)
	_, _ = oss.Init(global.App.Config.Storage.Disks.AliOss)
}
