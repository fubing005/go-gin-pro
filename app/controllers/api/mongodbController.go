package api

/*
	文档：
	https://mongoing.com/archives/27257
	https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
	https://cloud.tencent.com/developer/article/1739231【测试】
*/

import (
	"context"
	"fmt"
	"shalabing-gin/app/common/response"
	"shalabing-gin/global"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBController struct {
}

// 定义插入数据的结构体
type sunshareboy struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

var ctx = context.TODO()

func (con MongoDBController) GetDocuments(c *gin.Context) {
	// result := sunshareboy{
	// 	Name: "",
	// 	Age:  0,
	// 	City: "",
	// }

	//查询单个文档
	// filter := bson.D{
	// 	{"name", "wanger"},
	// }
	// err := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare").FindOne(context.TODO(), filter).Decode(&result)

	//查询多个文档
	//定义返回文档数量
	findOptions := options.Find()
	findOptions.SetLimit(5)
	//定义一个切片存储结果
	var results []*sunshareboy
	//将bson.D{{}}作为一个filter来匹配所有文档
	cur, err := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	//查找多个文档返回一个游标
	//遍历游标一次解码一个游标
	for cur.Next(context.TODO()) {
		//定义一个文档，将单个文档解码为result
		result := sunshareboy{
			Name: "",
			Age:  0,
			City: "",
		}
		err := cur.Decode(&result)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		results = append(results, &result)
	}
	if err := cur.Err(); err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	//遍历结束后关闭游标
	cur.Close(context.TODO())

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, results)
}

func (con MongoDBController) AddDocument(c *gin.Context) {

	//连接到test库的sunshare集合,集合不存在会自动创建
	collection := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare")

	// wanger := sunshareboy{"wanger", 24, "北京"}
	// insertOne, err := collection.InsertOne(ctx, wanger)
	dongdong := sunshareboy{"张冬冬", 29, "成都"}
	huazai := sunshareboy{"华仔", 28, "深圳"}
	suxin := sunshareboy{"素心", 24, "甘肃"}
	god := sunshareboy{"刘大仙", 24, "杭州"}
	qiaoke := sunshareboy{"乔克", 29, "重庆"}
	jiang := sunshareboy{"姜总", 24, "上海"}

	//插入多条数据要用到切片
	boys := []interface{}{dongdong, huazai, suxin, god, qiaoke, jiang}
	insertMany, err := collection.InsertMany(ctx, boys)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	// response.Success(c, insertOne.InsertedID)
	response.Success(c, insertMany.InsertedIDs)
}

func (con MongoDBController) UpdateDocument(c *gin.Context) {
	//更新单个文档
	// filter := bson.D{{Key: "name", Value: "张冬冬"}}
	// //如果过滤的文档不存在，则插入新的文档
	// opts := options.Update().SetUpsert(true)
	// update := bson.D{
	// 	{Key: "$set", Value: bson.D{
	// 		{Key: "city", Value: "北京"}},
	// 	}}
	// result, err := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare").UpdateOne(context.TODO(), filter, update, opts)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	// res1 := ""
	// res2 := ""
	// if result.MatchedCount != 0 {
	// 	res1 = fmt.Sprintf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	// }
	// if result.UpsertedCount != 0 {
	// 	res2 = fmt.Sprintf("inserted a new document with ID %v\n", result.UpsertedID)
	// }

	// response.Success(c, gin.H{
	// 	"res1": res1,
	// 	"res2": res2,
	// })

	//更新多个文档
	filter := bson.D{{Key: "city", Value: "北京"}}
	//如果过滤的文档不存在，则插入新的文档
	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "city", Value: "铁岭"}},
		}}
	result, err := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare").UpdateMany(context.TODO(), filter, update, opts)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	res1 := ""
	res2 := ""
	if result.MatchedCount != 0 {
		res1 = fmt.Sprintf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	}
	if result.UpsertedCount != 0 {
		res2 = fmt.Sprintf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	response.Success(c, gin.H{
		"res1": res1,
		"res2": res2,
	})
}

func (con MongoDBController) DeleteDocument(c *gin.Context) {
	filter := bson.D{{Key: "city", Value: "铁岭"}}
	deleteResult, err := global.App.MongoDB.Database("shalabing_gin").Collection("sunshare").DeleteMany(context.TODO(), filter)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	res := fmt.Sprintf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	response.Success(c, res)
}

func (con MongoDBController) GetServiceStatus(c *gin.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	// 使用 bson.D 构建命令
	command := bson.D{
		{Key: "serverStatus", Value: 1},
	}
	// 执行命令
	singleResult := global.App.MongoDB.Database("admin").RunCommand(ctx, command)
	err := singleResult.Err()
	if err != nil {
		res1 := fmt.Sprint("Failed to run serverStatus command: %v", err)
		response.BusinessFail(c, res1)
		return
	}

	var serverStatus bson.M
	// 解码结果
	err = singleResult.Decode(&serverStatus)
	if err != nil {
		res1 := fmt.Sprintf("Failed to decode serverStatus result: %v", err)
		response.BusinessFail(c, res1)
		return
	}

	version, ok := serverStatus["version"].(string)
	if !ok {
		response.BusinessFail(c, "Failed to get MongoDB version from server status")
		return
	}
	response.Success(c, gin.H{
		"serverStatus": serverStatus,
		"version":      version,
	})
}
