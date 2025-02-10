package bootstrap

import (
	"shalabing-gin/global"

	"github.com/elastic/go-elasticsearch/v7"
	"go.uber.org/zap"
)

func InitializeES() *elasticsearch.Client {
	// 初始化 Elasticsearch 客户端
	esConfig := elasticsearch.Config{
		Addresses: global.App.Config.Elasticsearch.Addresses,
		Username:  global.App.Config.Elasticsearch.Username,
		Password:  global.App.Config.Elasticsearch.Password,
		// 配置HTTP传输对象
		// Transport: &http.Transport{
		// 	//MaxIdleConnsPerHost 如果非零，控制每个主机保持的最大空闲(keep-alive)连接。如果为零，则使用默认配置2。
		// 	MaxIdleConnsPerHost: 10,
		// 	//ResponseHeaderTimeout 如果非零，则指定在写完请求(包括请求体，如果有)后等待服务器响应头的时间。
		// 	ResponseHeaderTimeout: time.Second,
		// 	//DialContext 指定拨号功能，用于创建不加密的TCP连接。如果DialContext为nil(下面已弃用的Dial也为nil)，那么传输拨号使用包网络。
		// 	DialContext: (&net.Dialer{Timeout: time.Second}).DialContext,
		// 	// // TLSClientConfig指定TLS.client使用的TLS配置。
		// 	// //如果为空，则使用默认配置。
		// 	// //如果非nil，默认情况下可能不启用HTTP/2支持。
		// 	// TLSClientConfig: &tls.Config{
		// 	// 	MaxVersion: tls.VersionTLS11,
		// 	// 	//InsecureSkipVerify 控制客户端是否验证服务器的证书链和主机名。
		// 	// 	InsecureSkipVerify: true,
		// 	// },
		// },
	}

	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		global.App.Log.Error("Elasticsearch 客户端初始化失败:", zap.Any("err", err))
	}

	// 测试连接
	// res, err := es.Info()
	// if err != nil {
	// 	global.App.Log.Error("无法连接到 Elasticsearch: ", zap.Any("err", err))
	// }
	// if res == nil {
	// 	global.App.Log.Error("Elasticsearch 响应为空，请检查服务是否正常启动")
	// }
	// // 确保响应的 Body 被关闭
	// defer func() {
	// 	if res.Body != nil {
	// 		res.Body.Close()
	// 	}
	// }()
	global.App.Log.Info("Elasticsearch 连接成功")
	return es
}
