package utils

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"shalabing-gin/app/models"
	"shalabing-gin/global"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/ini.v1"

	//引入模块的时候前面加个.表示可以直接使用模块里面的方法，无需加模块名称
	"github.com/gomarkdown/markdown"
	. "github.com/hunterhug/go_image"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	errors_tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

// 时间戳转换成日期函数
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 日期转换成时间戳
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 获取当前时间戳(毫秒)
func GetUnix() int64 {
	return time.Now().Unix()
}

// 获取当前时间戳(纳秒)
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// md5加密
func Md5(str string) string {
	//data := []byte(str)
	//return fmt.Sprintf("%x\n", md5.Sum(data))

	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 表示把string字符串转换成int
func Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// 表示把string字符串转换成Float
func Float(str string) (float64, error) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err
}

// 表示把int转换成string字符串
func String(n int) string {
	str := strconv.Itoa(n)
	return str
}

// 图片上传:上传到cos
func CosUploadImg(c *gin.Context, picName string) (string, error) {
	//1.获取上传文件
	file, err := c.FormFile(picName)
	//判断上传文件上否存在
	if err != nil { //说明上传文件不存在
		return "", nil
	}
	//2.获取后缀名,判断后缀是否正确: .jpg,.png,.gif,.jpeg
	extName := path.Ext(file.Filename)
	//设置后缀map
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	//判断后缀是否合法
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}
	//3.创建图片保存目录 ./static/upload/20230203
	//获取日期
	day := GetDay()
	//拼接目录, 上传时,cos会自动创建对应文件目录
	dir := "./static/upload/" + day
	//4.生成文件名称和文件保存目录: models.GetUnixNano() 获取时间戳(int64) 纳秒:防止速度过快而上传图片失败; strconv.FormatInt() 把时间戳(int64)转换成字符串
	filename := strconv.FormatInt(GetUnixNano(), 10) + extName
	//5.执行上传
	dst := path.Join(dir, filename)
	//上传文件到指定目录
	_, err1 := CosUpload(file, dst)
	if err1 != nil {
		return "", err1
	}
	fmt.Println(dst)
	return dst, nil
}

// 图片上传:上传到OSS
func OssUploadImg(c *gin.Context, picName string) (string, error) {
	//1.获取上传文件
	file, err := c.FormFile(picName)
	//判断上传文件上否存在
	if err != nil { //说明上传文件不存在
		return "", nil
	}
	//2.获取后缀名,判断后缀是否正确: .jpg,.png,.gif,.jpeg
	extName := path.Ext(file.Filename)
	//设置后缀map
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	//判断后缀是否合法
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}
	//3.创建图片保存目录 ./static/upload/20230203
	//获取日期
	day := GetDay()
	//拼接目录, 上传时,oss会自动创建对应文件目录
	dir := "./static/upload/" + day
	//4.生成文件名称和文件保存目录: models.GetUnixNano() 获取时间戳(int64) 纳秒:防止速度过快而上传图片失败; strconv.FormatInt() 把时间戳(int64)转换成字符串
	filename := strconv.FormatInt(GetUnixNano(), 10) + extName
	//5.执行上传
	dst := path.Join(dir, filename)

	//上传文件到指定目录
	OssUpload(file, dst)
	return dst, nil
}

// 封装oss上传图片方法
func OssUpload(file *multipart.FileHeader, dst string) (string, error) {
	// 1.创建OSSClient实例。
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "xxx", "xxx")
	if err != nil {
		return "", err
	}
	// 2.获取存储空间。
	bucket, err := client.Bucket("gin")
	if err != nil {
		return "", err
	}

	// 3.读取本地文件: file.Open()返回的File最终的类型为:io.Reader， 这样下面的上传就可以用了
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 上传文件流 bucket.PutObjec方法第二个参数类型为io.Reader, src的类型为
	err = bucket.PutObject(dst, src)
	if err != nil {
		return "", err
	}
	return dst, nil
}

// 封装cos上传图片方法
func CosUpload(file *multipart.FileHeader, dst string) (string, error) {
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	// 1.创建CosClient实例。
	//获取cos相关配置
	cfg, err := ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	rawURL := cfg.Section("cos").Key("rawURL").String()
	secretID := cfg.Section("cos").Key("secretID").String()
	scretKey := cfg.Section("cos").Key("scretKey").String()
	u, _ := url.Parse(rawURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，
			//登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: secretID, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，
			//登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: scretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})

	//对象在存储桶中的唯一标识
	key := dst
	//通过文件流上传对象
	// 3.读取本地文件: file.Open()返回的File最终的类型为:io.Reader， 这样下面的上传就可以用了
	fd, errOpen := file.Open()
	if errOpen != nil {
		return "", errOpen
	}
	defer fd.Close()

	_, err = client.Object.Put(context.Background(), key, fd, nil)
	if err != nil {
		return "", err
	}

	return dst, nil
}

// 图片上传:上传到本地
func LocalUploadImg(c *gin.Context, picName string) (string, error) {
	//1.获取上传文件
	file, err := c.FormFile(picName)
	//判断上传文件上否存在
	if err != nil { //说明上传文件不存在
		return "", nil
	}
	//2.获取后缀名,判断后缀是否正确: .jpg,.png,.gif,.jpeg
	extName := path.Ext(file.Filename)
	//设置后缀map
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	//判断后缀是否合法
	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}
	//3.创建图片保存目录 ./static/upload/20230203
	//获取日期
	day := GetDay()
	//拼接目录
	dir := "./static/upload/" + day
	//创建目录:MkdirAll 目录不存在,会一次性创建多层
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		return "", err
	}
	//4.生成文件名称和文件保存目录: models.GetUnixNano() 获取时间戳(int64) 纳秒:防止速度过快而上传图片失败; strconv.FormatInt() 把时间戳(int64)转换成字符串
	filename := strconv.FormatInt(GetUnixNano(), 10) + extName
	//5.执行上传
	dst := path.Join(dir, filename)
	//上传文件到指定目录
	c.SaveUploadedFile(file, dst)
	return dst, nil
}

// 通过列获取系统设置里面的值,columnName就是结构体的属性名称
func GetSettingFromColumn(columnName string) string {
	setting := models.Setting{}
	global.App.DB.First(&setting)
	//反射来获取
	v := reflect.ValueOf(setting)
	val := v.FieldByName(columnName).String()
	return val
}

// 生成商品缩略图
func ResizeGoodsImage(filename string) {
	//获取文件后缀名
	extname := path.Ext(filename)
	//获取缩略图尺寸
	thumbnailSizeSlice := strings.Split(GetSettingFromColumn("ThumbnailSize"), ",")
	//static/upload/tao_400.png
	//static/upload/tao_400.png_100x100.png
	//遍历尺寸,生成缩略图
	for i := 0; i < len(thumbnailSizeSlice); i++ {
		savepath := filename + "_" + thumbnailSizeSlice[i] + "x" + thumbnailSizeSlice[i] + extname
		w, _ := Int(thumbnailSizeSlice[i])
		//调用github.com/hunterhug/go_image中的方法ThumbnailF2F(),生成缩略图
		err := ThumbnailF2F(filename, savepath, w, w)
		if err != nil {
			fmt.Println(err) //写个日志模块  处理日志
		}
	}
}

// SubStr截取字符串
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	rl := len(rs)
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = 0
	}

	if end < 0 {
		end = rl
	}
	if end > rl {
		end = rl
	}
	if start > end {
		start, end = end, start
	}
	return string(rs[start:end])
}

/*
转换markdown格式
str就是markdown语法
例如:
**我是一个三级标题**  <=> <h3>我是一个三级标题</h3>
**我是一个加粗** <=> <strong>我是一个加粗</strong>
*/
func FormatAttr(str string) string {
	tempSlice := strings.Split(str, "\n")
	var tempStr string
	for _, v := range tempSlice {
		md := []byte(v)
		output := markdown.ToHTML(md, nil, nil)
		tempStr += string(output)
	}
	return tempStr
}

// 发送短信:腾讯云短信
func SmsTencent(phone string, code string) (bool, error) {
	cfg, err := ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	secretId := cfg.Section("sms_tencent").Key("secretId").String()
	secretKey := cfg.Section("sms_tencent").Key("secretKey").String()
	endpoint := cfg.Section("sms_tencent").Key("endpoint").String()
	signMethod := cfg.Section("sms_tencent").Key("signMethod").String()
	region := cfg.Section("sms_tencent").Key("region").String()
	smsSdkAppId := cfg.Section("sms_tencent").Key("smsSdkAppId").String()
	signName := cfg.Section("sms_tencent").Key("signName").String()
	templateId_common := cfg.Section("sms_tencent").Key("templateId_common").String()

	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
	 * 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	 * 以免泄露密钥对危及你的财产安全。
	 * SecretId、SecretKey 查询: https://console.cloud.tencent.com/cam/capi */
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5

	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = endpoint

	/* SDK默认用TC3-HMAC-SHA256进行签名
	 * 非必要请不要修改这个字段 */
	cpf.SignMethod = signMethod

	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := sms.NewClient(credential, region, cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 你可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/smsv2
	 * 腾讯云短信小助手: https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81 */

	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	// 应用 ID 可前往 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage) 查看
	request.SmsSdkAppId = common.StringPtr(smsSdkAppId)

	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	// 签名信息可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-sign) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-sign) 的签名管理查看
	request.SignName = common.StringPtr(signName)

	/* 模板 ID: 必须填写已审核通过的模板 ID */
	// 模板 ID 可前往 [国内短信](https://console.cloud.tencent.com/smsv2/csms-template) 或 [国际/港澳台短信](https://console.cloud.tencent.com/smsv2/isms-template) 的正文模板管理查看
	request.TemplateId = common.StringPtr(templateId_common)

	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs([]string{code})

	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})

	/* 用户的 session 内容（无需要可忽略）: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	request.SessionContext = common.StringPtr("")

	/* 短信码号扩展号（无需要可忽略）: 默认未开通，如需开通请联系 [腾讯云短信小助手] */
	request.ExtendCode = common.StringPtr("")

	/* 国内短信无需填写该项；国际/港澳台短信已申请独立 SenderId 需要填写该字段，默认使用公共 SenderId，无需填写该字段。注：月度使用量达到指定量级可申请独立 SenderId 使用，详情请联系 [腾讯云短信小助手](https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81)。 */
	request.SenderId = common.StringPtr("")

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors_tencent.TencentCloudSDKError); ok {
		return false, err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		return false, err
	}
	//返回结果
	result := response.Response
	// 当前只有一条
	statusSet := result.SendStatusSet
	//返回的结果
	result_code := *statusSet[0].Code
	if result_code != "Ok" {
		return false, errors.New(*statusSet[0].Message)
	}

	return true, nil
	/* 当出现以下错误码时，快速解决方案参考
	 * [FailedOperation.SignatureIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.signatureincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [FailedOperation.TemplateIncorrectOrUnapproved](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Afailedoperation.templateincorrectorunapproved-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnauthorizedOperation.SmsSdkAppIdVerifyFail](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunauthorizedoperation.smssdkappidverifyfail-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * [UnsupportedOperation.ContainDomesticAndInternationalPhoneNumber](https://cloud.tencent.com/document/product/382/9558#.E7.9F.AD.E4.BF.A1.E5.8F.91.E9.80.81.E6.8F.90.E7.A4.BA.EF.BC.9Aunsupportedoperation.containdomesticandinternationalphonenumber-.E5.A6.82.E4.BD.95.E5.A4.84.E7.90.86.EF.BC.9F)
	 * 更多错误，可咨询[腾讯云助手](https://tccc.qcloud.com/web/im/index.html#/chat?webAppId=8fa15978f85cb41f7e2ea36920cb3ae1&title=Sms)
	 */
}

// 生成随机数
func GetRandomNum() string {
	var str string
	for i := 0; i < 4; i++ {
		current := rand.Intn(10)
		str += strconv.Itoa(current)
	}
	return str
}

// 获取订单编号
func GetOrderId() string {
	// 2022020312233
	temp := "20060102150405"
	return time.Now().Format(temp) + GetRandomNum()
}

func StructToJsonString(s interface{}) string {
	data, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(data)
}

func Contains(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
