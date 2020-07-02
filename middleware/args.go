package middleware

import (
	"encoding/json"
	"github.com/shiningacg/mygin"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	ArgsPrefix = "ARGS_"
)

// 负责参数校验
func Args() mygin.HandlerFunc {
	return func(context *mygin.Context) {
		if context.Request.Method == "POST" {
			args, err := parseBodyJson(context.Request.Body)
			if err != nil {
				log.Println(err)
				context.Abort()
			}
			for key, value := range args {
				context.Set(ArgsPrefix+key, value)
			}
			context.Next()
		}
	}
}

func key(k string) string {
	return ArgsPrefix + k
}

// TODO：添加tag类验证支持
func Merge(ctx *mygin.Context, target interface{}) error {
	v := reflect.ValueOf(target).Elem()
	t := reflect.TypeOf(target).Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		// 判断filed的类型是不是指针
		vl := v.Field(i)
		tp := t.Field(i)
		var value interface{}
		// 字段名称
		Name := tp.Name
		// Tag名
		Tag := tp.Tag.Get("json")
		// 寻找字段
		if temp := ctx.Value(key(strings.ToLower(Name))); temp != nil {
			value = temp
		} else if temp := ctx.Value(key(Tag)); temp != nil {
			value = temp
		} else {
			continue
		}
		val, ok1 := value.(string)
		val2, ok2 := value.(float64)
		// 如果都不成功，那么就是结构体类型
		if !ok1 && !ok2 {
			_v := reflect.ValueOf(value)
			if _v.Kind() == reflect.Ptr && _v.Type() == tp.Type {
				vl.Set(_v)
			}
			continue
		}
		// 如果是字符串类型，那么进行可能的类型转化
		if ok1 {
			switch vl.Kind() {
			case reflect.Int:
				temp, err := strconv.ParseInt(val, 0, 64)
				if err != nil {
					break
				}
				v.Field(i).SetInt(temp)
			case reflect.String:
				v.Field(i).SetString(val)
			case reflect.Float64:
				temp, err := strconv.ParseFloat(val, 64)
				if err != nil {
					break
				}
				v.Field(i).SetFloat(temp)
			}
		}
		// 如果是float64类型
		if ok2 {
			switch vl.Kind() {
			case reflect.Int:
				v.Field(i).SetInt(int64(val2))
			case reflect.Float64:
				v.Field(i).SetFloat(val2)
			}
		}
	}
	return nil
}

func readAll(dst []byte, src io.Reader) (int, error) {
	var err error
	var total, n int
	for {
		if cap(dst[total:]) == 0 {
			return 0, mygin.ErrReachLimitSize
		}
		n, err = src.Read(dst[total:])
		if err != nil {
			if err == io.EOF {
				total += n
				break
			}
			return 0, err
		}
		total += n
	}
	return total, nil
}

func parseBodyJson(reader io.Reader) (map[string]interface{}, error) {
	var err error
	var n int
	temp := make([]byte, 1024*100)
	n, err = readAll(temp, reader)
	if err != nil {
		return nil, err
	}
	// 拿到数据
	data := temp[:n]
	args := make(map[string]interface{})
	err = json.Unmarshal(data, &args)
	if err != nil {
		return nil, err
	}
	return args, nil
}
