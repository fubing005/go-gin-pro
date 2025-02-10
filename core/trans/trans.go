package trans

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"shalabing-gin/global"
	"strings"
	"sync"
)

var (
	translationsCache = make(map[string]map[string]interface{})
	cacheMutex        sync.RWMutex
)

// Trans 翻译
func Trans(key string, values ...interface{}) string {
	lang := global.App.Config.App.Lang
	// 使用读锁读取缓存
	cacheMutex.RLock()
	translations, exists := translationsCache[lang]
	cacheMutex.RUnlock()

	if !exists {
		// 缓存中没有，则加载文件并解析
		path := buildFilePath(lang)
		translations = loadTranslations(path)
		// 写锁更新缓存
		cacheMutex.Lock()
		translationsCache[lang] = translations
		cacheMutex.Unlock()
	}
	// 查找key对应的值
	// 按点分割键
	keys := strings.Split(key, ".")
	// 逐级查找翻译值
	for _, k := range keys {
		value, exists := translations[k]
		if !exists {
			return ""
		}

		// 如果是 map 类型，继续查找
		if nested, ok := value.(map[string]interface{}); ok {
			translations = nested
		} else {
			// 如果不是 map 类型，说明找到了最终值
			if str, ok := value.(string); ok {
				return fmt.Sprintf(str, values...)
			}
			fmt.Printf("获取翻译出错！%v", values)
		}
	}
	return ""
}

// buildFilePath 构建文件路径
func buildFilePath(lang string) string {
	var builder strings.Builder
	builder.WriteString("lang\\")
	builder.WriteString(lang)
	builder.WriteString("\\")
	builder.WriteString(lang)
	builder.WriteString(".json")
	return builder.String()
}

// loadTranslations 加载并解析翻译文件
func loadTranslations(path string) map[string]interface{} {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件失败:", err.Error())
		return nil
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("读取文件内容失败:", err.Error())
		return nil
	}

	// 解析JSON
	var translations map[string]interface{}
	if err := json.Unmarshal(bytes, &translations); err != nil {
		fmt.Println("解析JSON失败:", err.Error())
		return nil
	}

	return translations
}
