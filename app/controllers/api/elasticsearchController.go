package api

/*
	https://golang.halfiisland.com/community/database/Elasticsearch.html
*/

import (
	"bytes"
	"encoding/json"
	"io"
	"shalabing-gin/app/common/response"
	"shalabing-gin/global"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ElasticsearchController struct{}

// 创建文档
type Book struct {
	ID      string     `json:"id"`
	Author  string     `json:"author"`
	Name    string     `json:"name"`
	Pages   int        `json:"pages"`
	Price   float64    `json:"price"`
	PubDate *time.Time `json:"pubDate"`
	Summary string     `json:"summary"`
}

// 下面的代码借助go json的omitempty，在将更新数据对象序列化成json，
// 可以只序列化非零值字段，实现局部更新。 实际项目采用这种方式时，需要注意某个字段的零值具有业务意义时，可以采用对应的指针类型实现。
type doc struct {
	Doc interface{} `json:"doc"`
}
type Book1 struct {
	Author  string     `json:"author,omitempty"`
	Name    string     `json:"name,omitempty"`
	Pages   int        `json:"pages,omitempty"`
	Price   float64    `json:"price,omitempty"`
	PubDate *time.Time `json:"pubDate,omitempty"`
	Summary string     `json:"summary,omitempty"`
}

// cat
func (con ElasticsearchController) Cat(c *gin.Context) {
	//info
	info, _ := global.App.ES.Info(global.App.ES.Info.WithHuman())
	infoBody, _ := io.ReadAll(info.Body)
	infoData := map[string]interface{}{}
	_ = json.Unmarshal(infoBody, &infoData)

	// 查看集群健康状态
	health, _ := global.App.ES.Cat.Health()
	healthBody, _ := io.ReadAll(health.Body)

	//allocation 查看节点磁盘资源分配信息
	allocation, _ := global.App.ES.Cat.Allocation(global.App.ES.Cat.Allocation.WithNodeID(""))
	allocationBody, _ := io.ReadAll(allocation.Body)

	//shards 节点分片信息
	shards, _ := global.App.ES.Cat.Shards()
	shardsBody, _ := io.ReadAll(shards.Body)

	//master?v 查看主节点
	master, _ := global.App.ES.Cat.Master()
	masterBody, _ := io.ReadAll(master.Body)

	// nodes 查看集群信息
	nodes, _ := global.App.ES.Cat.Nodes()
	nodesBody, _ := io.ReadAll(nodes.Body)

	// tasks 返回集群中一个或多个节点上当前执行的任务信息。
	tasks, _ := global.App.ES.Cat.Tasks()
	tasksBody, _ := io.ReadAll(tasks.Body)

	// indices 查看索引信息
	indices, _ := global.App.ES.Cat.Indices(global.App.ES.Cat.Indices.WithIndex(""))
	indicesBody, _ := io.ReadAll(indices.Body)

	// segments 分片中的分段信息
	segments, _ := global.App.ES.Cat.Segments()
	segmentsBody, _ := io.ReadAll(segments.Body)

	// count 查看当前集群的doc数量
	count, _ := global.App.ES.Cat.Count()
	countBody, _ := io.ReadAll(count.Body)

	// recovery 显示正在进行和先前完成的索引碎片恢复的视图
	recovery, _ := global.App.ES.Cat.Recovery()
	recoveryBody, _ := io.ReadAll(recovery.Body)

	// /pending_tasks 显示正在等待的任务
	pendingTasks, _ := global.App.ES.Cat.PendingTasks()
	pendingTasksBody, _ := io.ReadAll(pendingTasks.Body)

	// aliases 显示别名、过滤器、路由信息
	aliases, _ := global.App.ES.Cat.Aliases()
	aliasesBody, _ := io.ReadAll(aliases.Body)

	// thread_pool 查看线程池信息
	threadPool, _ := global.App.ES.Cat.ThreadPool()
	threadPoolBody, _ := io.ReadAll(threadPool.Body)

	// plugins 显示每个运行插件节点的视图
	plugins, _ := global.App.ES.Cat.Plugins()
	pluginsBody, _ := io.ReadAll(plugins.Body)

	// fielddata #查看每个数据节点上fielddata当前占用的堆内存。
	fielddata, _ := global.App.ES.Cat.Fielddata()
	fielddataBody, _ := io.ReadAll(fielddata.Body)

	// nodeattrs 查看每个节点的自定义属性
	nodeattrs, _ := global.App.ES.Cat.Nodeattrs()
	nodeattrsBody, _ := io.ReadAll(nodeattrs.Body)

	// /templates 模板
	templates, _ := global.App.ES.Cat.Templates()
	templatesBody, _ := io.ReadAll(templates.Body)

	result := gin.H{
		"info":         infoData,
		"health":       string(healthBody),
		"allocation":   string(allocationBody),
		"shards":       string(shardsBody),
		"master":       string(masterBody),
		"nodes":        string(nodesBody),
		"tasks":        string(tasksBody),
		"indices":      string(indicesBody),
		"segments":     string(segmentsBody),
		"count":        string(countBody),
		"recovery":     string(recoveryBody),
		"pendingTasks": string(pendingTasksBody),
		"aliases":      string(aliasesBody),
		"threadPool":   string(threadPoolBody),
		"plugins":      string(pluginsBody),
		"fielddata":    string(fielddataBody),
		"nodeattrs":    string(nodeattrsBody),
		"templates":    string(templatesBody),
	}

	response.Success(c, result)
}

// 查询索引
func (con ElasticsearchController) GetIndex(c *gin.Context) {
	// type IndicesGet func(index []string, o ...func(*IndicesGetRequest)) (*Response, error)

	//Get 获取指定索引的映射
	result, _ := global.App.ES.API.Indices.Get([]string{"es_test_book"})
	body, _ := io.ReadAll(result.Body)
	jsonData := map[string]interface{}{}
	json.Unmarshal(body, &jsonData)
	response.Success(c, jsonData)
}

// 创建索引
func (con ElasticsearchController) CreateIndex(c *gin.Context) {
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type": "keyword",
				},
				"pages": map[string]interface{}{
					"type": "integer",
				},
				"price": map[string]interface{}{
					"type": "float",
				},
				"pubDate": map[string]interface{}{
					"type": "date",
				},
				"summary": map[string]interface{}{
					"type":            "text",
					"analyzer":        "ik_max_word",
					"search_analyzer": "ik_smart",
				},
			},
		},
	})

	result, _ := global.App.ES.API.Indices.Create("es_test_book", global.App.ES.API.Indices.Create.WithBody(body))

	response.Success(c, result)
}

// 删除索引
func (con ElasticsearchController) DeleteIndex(c *gin.Context) {
	result, _ := global.App.ES.API.Indices.Delete([]string{"poet-index"})
	response.Success(c, result)
}

// 查询文档
func (con ElasticsearchController) GetDocuments(c *gin.Context) {
	//1. 查询文档
	id := c.Query("id")
	result, err := global.App.ES.Get("es_test_book", id)
	body, _ := io.ReadAll(result.Body)
	jsonData := map[string]interface{}{}
	json.Unmarshal(body, &jsonData)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	//2.文档搜索 单个字段
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"name": "神雕侠侣",
	// 		},
	// 	},
	// 	"highlight": map[string]interface{}{
	// 		"pre_tags":  []string{"<font color='red'>"},
	// 		"post_tags": []string{"</font>"},
	// 		"fields": map[string]interface{}{
	// 			"name": map[string]interface{}{},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//3. 提高某个字段权重，可以使用 ^ 字符语法为单个字段提升权重，在字段名称的末尾添加 ^boost ，其中 boost 是一个浮点数：
	// /*
	// 	当搜索结果包含大量文档时，合理地提升某些字段的权重可以让搜索结果的排序更加符合用户的预期。比如在一个商品搜索系统中，商品的品牌和价格是用户比较关注的因素。如果提升品牌字段的权重，那么与搜索关键词匹配的知名品牌商品会更靠前显示，即使其他商品在其他字段上也有一定的匹配度，但由于品牌字段权重较高，知名品牌商品的整体相关性得分会更高，从而在搜索结果中排名更靠前
	// */
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":  "哈利・波特",
	// 			"fields": []string{"name", "author^2"},
	// 		},
	// 	},
	// 	"highlight": map[string]interface{}{
	// 		"pre_tags":  []string{"<font color='red'>"},
	// 		"post_tags": []string{"</font>"},
	// 		"fields": map[string]interface{}{
	// 			"name": map[string]interface{}{},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//4. 显示所有文档
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//5. 范围查询通过range实现范围查询，类似SQL语句中的>, >=, <, <=表达式。 gte范围参数 - 等价于>= lte范围参数 - 等价于 <= 范围参数可以只写一个，例如：仅保留 “gte”: 10， 则代表 FIELD字段 >= 10
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"range": map[string]interface{}{
	// 			"price": map[string]interface{}{
	// 				"gte": "0",
	// 				"lte": "97",
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//6. 文档搜索 多个字段
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":  "金庸", //eq 查询
	// 			"fields": []string{"name", "author"},
	// 		},
	// 	},
	// 	"highlight": map[string]interface{}{
	// 		"pre_tags":  []string{"<font color='red'>"},
	// 		"post_tags": []string{"</font>"},
	// 		"fields": map[string]interface{}{
	// 			"name": map[string]interface{}{},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//7. 通过term查询【不分词，精确匹配】,mathch是分词匹配
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"term": map[string]interface{}{
	// 			"author": "金庸",
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//8. 通过match查询
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"name": "《天龙八部》第二部",
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//9. 通过match_phrase查询
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match_phrase": map[string]interface{}{
	// 			"summary": map[string]interface{}{
	// 				"query": "末年南宋",
	// 				"slop":  2, //slop < 2 ,要求分词顺序是一致的，比如【末年南宋】就查不出数据，slop >= 2，可以忽略分词顺序
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//10. 通过multi_match查询
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"multi_match": map[string]interface{}{
	// 			"query":  "射雕英雄传",
	// 			"fields": []string{"name", "author"},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//11. 通过range查询
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"range": map[string]interface{}{
	// 			"price": map[string]interface{}{
	// 				"gte": 50,
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//12. 通过exists查询: 该字段不为 null 且不是空数组或空字符串的文档
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"exists": map[string]interface{}{
	// 			"field": "price",
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//13.布尔 must and
	// name := c.Query("name")
	// price := c.Query("price")
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"must": []map[string]interface{}{
	// 				{
	// 					"match": map[string]interface{}{
	// 						"author": name,
	// 					},
	// 				},
	// 				{
	// 					"match": map[string]interface{}{
	// 						"price": price,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{"pubDate": []map[string]interface{}{
	// 			{"order": "desc"},
	// 		},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//14.布尔 should or
	// name := c.Query("name")
	// price := c.Query("price")
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"should": []map[string]interface{}{
	// 				{
	// 					"match": map[string]interface{}{
	// 						"author": name,
	// 					},
	// 				},
	// 				{
	// 					"match": map[string]interface{}{
	// 						"price": price,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"pubDate": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//15.布尔 must not: 排除 name="金庸" and price=20
	// name := "金庸"
	// price := "20"
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"must_not": []map[string]interface{}{
	// 				{
	// 					"match": map[string]interface{}{
	// 						"author": name,
	// 					},
	// 				},
	// 				{
	// 					"match": map[string]interface{}{
	// 						"price": price,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"pubDate": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//16.精确度查询
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"term": map[string]interface{}{
	// 			"name": "射雕英雄传",
	// 		},
	// 	},
	// }

	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//17.使用Mysql 的方式查询
	// index := "es_test_book"
	// query := map[string]interface{}{
	// 	"query": "select count(1) from " + index,
	// }
	// jsonBody, _ := json.Marshal(query)
	// result, err := global.App.ES.SQL.Query(bytes.NewReader(jsonBody), global.App.ES.SQL.Query.WithContext(context.Background()))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//18. 聚合查询 - 聚合 group by
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_author": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"terms": map[string]interface{}{ // 聚合类型为，terms，terms是桶聚合的一种，类似SQL的group by的作用，根据字段分组，相同字段值的文档分为一组。
	// 				"field": "author", // terms聚合类型的参数，这里需要设置分组的字段为author，根据author分组
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"pubDate": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//19. 聚合查询 - 聚合 count
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_author": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"value_count": map[string]interface{}{
	// 				"field": "author",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{"pubDate": []map[string]interface{}{
	// 			{"order": "desc"},
	// 		},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// value, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(value, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//20. 聚合查询 - Cardinality 基数聚合，也是用于统计文档的总数，跟Value Count的区别是，基数聚合会去重，不会统计重复的值，类似SQL的count(DISTINCT 字段)用法。
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_author": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"cardinality": map[string]interface{}{
	// 				"field": "author", //
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 10,
	// 	"sort": []map[string]interface{}{
	// 		{"pubDate": []map[string]interface{}{
	// 			{"order": "desc"},
	// 		},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// value, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(value, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//21. 聚合查询 - Top Hits 最大值聚合 desc
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"aggregation_name": map[string]interface{}{
	// 			"terms": map[string]interface{}{
	// 				"field": "author",
	// 			},
	// 		},
	// 		"aggs": map[string]interface{}{
	// 			// "top_author": map[string]interface{}{
	// 			"top_hits": map[string]interface{}{
	// 				"size": 3,
	// 				"sort": []map[string]interface{}{
	// 					{
	// 						"price": []map[string]interface{}{
	// 							{"order": "desc"},
	// 						},
	// 					},
	// 				},
	// 				// },
	// 			},
	// 		},
	// 		// "from": 0,
	// 		// "size": 10,
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// value, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(value, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//22. 聚合查询 - Top Hits 最小值聚合 asc
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_author": map[string]interface{}{ // 给聚合查询取个名字）
	// 			"top_hits": map[string]interface{}{
	// 				"size": 1,
	// 				"sort": []map[string]interface{}{
	// 					{
	// 						"price": []map[string]interface{}{
	// 							{"order": "asc"},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		// "from": 0,
	// 		// "size": 1,
	// 		// "sort": []map[string]interface{}{
	// 		// 	{
	// 		// 		"price": []map[string]interface{}{
	// 		// 			{"order": "desc"},
	// 		// 		},
	// 		// 	},
	// 		// },
	// 	},
	// }

	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//23. 聚合查询 - Avg
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"avg": map[string]interface{}{
	// 				"field": "price", //
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//24. 聚合查询 - sum
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"sum": map[string]interface{}{
	// 				"field": "price",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//25. 聚合查询 - max
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"max_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"max": map[string]interface{}{
	// 				"field": "price", //
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"pubDate": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//26. 聚合查询 - min
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"min_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"min": map[string]interface{}{
	// 				"field": "price",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//27. 聚合查询 - 综合
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"min_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"min": map[string]interface{}{
	// 				"field": "price",
	// 			},
	// 		},
	// 		"count_name": map[string]interface{}{
	// 			"value_count": map[string]interface{}{
	// 				"field": "name",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 2,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "asc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//28. 聚合查询 - Terms聚合 terms聚合的作用跟SQL中group by作用一样，都是根据字段唯一值对数据进行分组（分桶），字段值相等的文档都分到同一个桶内。select price, count(*) from book group by price
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_author": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"terms": map[string]interface{}{
	// 				"field": "author",
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//29. 聚合查询 - Histogram直方图聚合[https://www.cnblogs.com/xing901022/p/4954823.html]
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"prices": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"histogram": map[string]interface{}{
	// 				"field":    "price",
	// 				"interval": 50, // 分桶的间隔为50，意思就是price字段值按50间隔分组
	// 				"min_doc_count" : 1, //如果不想要显示doc_count为0的桶，可以通过min_doc_count来设置
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//30. 聚合查询 - Range聚合
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"range_name": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"range": map[string]interface{}{
	// 				"field": "price",
	// 				"ranges": []map[string]interface{}{
	// 					{"to": 20},
	// 					{"from": 20, "to": 50},
	// 					{"from": 50},
	// 				},
	// 			},
	// 			"aggs": map[string]interface{}{
	// 				"avg_price": map[string]interface{}{
	// 					"avg": map[string]interface{}{
	// 						"field": "price",
	// 					},
	// 				},
	// 				"max_price": map[string]interface{}{
	// 					"max": map[string]interface{}{
	// 						"field": "price",
	// 					},
	// 				},
	// 				"min_price": map[string]interface{}{
	// 					"min": map[string]interface{}{
	// 						"field": "price",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//31. 聚合查询 - Date histogram
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_name": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"date_histogram": map[string]interface{}{
	// 				"field":             "pubDate",             // 根据date字段分组
	// 				"calendar_interval": "month",               // 分组间隔：month代表每月、支持minute（每分钟）、hour（每小时）、day（每天）、week（每周）、year（每年）
	// 				"format":            "yyyy-MM-dd HH:mm:ss", // 设置返回结果中桶key的时间格式
	// 				"time_zone":         "Asia/Shanghai",       // 设置时区
	// 				"min_doc_count":     1,                     // 设置最小文档数量
	// 			},
	// 			"aggs": map[string]interface{}{
	// 				"total_price": map[string]interface{}{
	// 					"sum": map[string]interface{}{
	// 						"field": "price",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//32. 综合案例
	// query := map[string]interface{}{
	// 	"aggs": map[string]interface{}{ // 合查询语句的简写
	// 		"count_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 			"range": map[string]interface{}{
	// 				"field": "price",
	// 				"ranges": []interface{}{ // 范围配置
	// 					map[string]float64{"to": 20.0},               // 意思就是 price <= 200的文档归类到一个桶
	// 					map[string]float64{"from": 20.0, "to": 50.0}, // price>20 and price<50的文档归类到一个桶
	// 					map[string]float64{"from": 200.0},            // price>200的文档归类到一个桶},
	// 				},
	// 			},
	// 			"aggs": map[string]interface{}{ // 合查询语句的简写
	// 				"min_price": map[string]interface{}{ // 给聚合查询取个名字，
	// 					"min": map[string]interface{}{
	// 						"field": "price",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"from": 0,
	// 	"size": 1,
	// 	"sort": []map[string]interface{}{
	// 		{
	// 			"price": []map[string]interface{}{
	// 				{"order": "desc"},
	// 			},
	// 		},
	// 	},
	// }
	// marshal, _ := json.Marshal(query)
	// result, err := global.App.ES.Search(global.App.ES.Search.WithIndex("es_test_book"), global.App.ES.Search.WithBody(bytes.NewReader(marshal)))
	// body, _ := io.ReadAll(result.Body)
	// jsonData := map[string]interface{}{}
	// json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	response.Success(c, jsonData)
}

// 创建文档
func (con ElasticsearchController) AddDocument(c *gin.Context) {
	// 1. 按照id创建文档
	body := &bytes.Buffer{}
	id := c.PostForm("id")
	author := c.PostForm("author")
	name := c.PostForm("name")
	summary := c.PostForm("summary")
	pubDate := time.Now()
	err := json.NewEncoder(body).Encode(&Book{
		Author:  author,
		Price:   96.0,
		Name:    name,
		Pages:   1978,
		PubDate: &pubDate,
		Summary: summary,
	})
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	// 创建
	result, err := global.App.ES.Create("es_test_book", id, body)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	//2.批量创建文档
	// /*
	// 	批量操作对应ES的REST API是:
	// 		https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html
	// 		// POST /<target>/_bulk
	// 		// { "index" : { "_id" : "1" } }
	// 		// { "field1" : "value1" }
	// 		// { "delete" : { "_id" : "2" } }
	// 		// { "create" : { "_id" : "3" } }
	// 		// { "field1" : "value3" }
	// 		// { "update" : {"_id" : "1" } }
	// 		// { "doc" : {"field2" : "value2"} }
	// 		// 对应index, create, update操作，
	// 		// 提交的数据都是由两行组成，第一行是meta数据，描述操作信息
	// 		// ，第二行是具体提交的数据，对于delete操作只有一行meta数据。 对照REST API
	// */
	// pubDate := time.Now()
	// createBooks := []*Book{
	// 	{
	// 		ID:      "6",
	// 		Name:    "神雕侠侣",
	// 		Author:  "金庸",
	// 		Pages:   1978,
	// 		Price:   96.0,
	// 		PubDate: &pubDate,
	// 		Summary: "杨过是杨康与穆念慈之子，杨康死后，穆念慈独自抚养杨过。穆念慈过世后，杨过被郭靖、黄蓉夫妇接到桃花岛。在桃花岛，杨过因性格叛逆，与武氏兄弟、郭芙等人产生冲突，且不受黄蓉信任，郭靖便将他送往终南山全真教学艺。然而，杨过在全真教饱受欺负，愤而逃入活死人墓，被小龙女收为徒弟。在墓中，杨过与小龙女朝夕相处，二人一同修炼古墓派武功。",
	// 	},
	// 	{
	// 		ID:      "7",
	// 		Name:    "射雕英雄传",
	// 		Author:  "金庸",
	// 		Pages:   1978,
	// 		Price:   96.0,
	// 		PubDate: &pubDate,
	// 		Summary: "南宋末年，临安府牛家村中，隐居于此的全真教道士丘处机与金国王爷完颜洪烈相遇，双方发生冲突。丘处机追杀完颜洪烈至嘉兴，与江南七怪在醉仙楼不打不相识。一番较量后，双方约定分别教导流落天涯的忠良之后郭啸天之子郭靖和杨铁心之子杨康，十八年后让两人在嘉兴再次比武，以分高下。",
	// 	},
	// }
	// body := &bytes.Buffer{}
	// for _, book := range createBooks {
	// 	meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, book.ID, "\n"))
	// 	data, err := json.Marshal(book)
	// 	if err != nil {
	// 		response.BusinessFail(c, err.Error())
	// 		return
	// 	}
	// 	data = append(data, "\n"...)
	// 	body.Grow(len(meta) + len(data))
	// 	body.Write(meta)
	// 	body.Write(data)
	// }
	// result, err := global.App.ES.Bulk(body, global.App.ES.Bulk.WithIndex("es_test_book"))
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	response.Success(c, result)
}

// 更新文档
func (con ElasticsearchController) UpdateDocument(c *gin.Context) {
	id := c.PostForm("id")
	// name := c.PostForm("name")
	price := c.PostForm("price")
	priceFloat, _ := strconv.ParseFloat(price, 64)
	//1. 覆盖性更新文档
	// body := &bytes.Buffer{}
	// author := c.PostForm("author")
	// pubDate := time.Now()
	// err := json.NewEncoder(body).Encode(&Book{
	// 	Author:  author,
	// 	Price:   97.0,
	// 	Name:    name,
	// 	Pages:   1978,
	// 	PubDate: &pubDate,
	// 	// Summary: "...",
	// })
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }
	// result, err := global.App.ES.Index("es_test_book", body, global.App.ES.Index.WithDocumentID(id))
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	//2.局部更新文档
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(&doc{
		Doc: &Book1{
			Price: priceFloat,
		},
	})
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	result, err := global.App.ES.Update("es_test_book", id, body)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, result)
}

// 删除文档
func (con ElasticsearchController) DeleteDocument(c *gin.Context) {
	//1. 删除文档
	result, err := global.App.ES.Delete("es_test_book", "1")
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	//2. 批量删除文档
	// body := &bytes.Buffer{}
	// deleteBookIds := []string{"1"}
	// for _, id := range deleteBookIds {
	// 	meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%s" } }%s`, id, "\n"))
	// 	body.Grow(len(meta))
	// 	body.Write(meta)
	// }
	// result, err := global.App.ES.Bulk(body, global.App.ES.Bulk.WithIndex("es_test_book"))
	// if err != nil {
	// 	response.BusinessFail(c, err.Error())
	// 	return
	// }

	response.Success(c, result)
}
