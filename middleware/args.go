package middleware

import (
	"encoding/json"
	"errors"
	"github.com/shiningacg/mygin"
	"io"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	ArgsPrefix = "ARGS_"
)

// 负责参数校验
func Args(target interface{}) mygin.HandlerFunc {
	return func(context *mygin.Context) {
		// 尝试从body内获取数据
		args, err := parseBodyJson(context.Request.Body)
		if err != nil {
			log.Println(err)
			context.Abort()
		}
		for k, v := range args {
			context.Set(key(k), v)
		}
		// 针对get方法的参数获取
		if context.Request.Method == "GET" {
			u, err := url.Parse(context.Request.RequestURI)
			if err != nil {
				return
			}
			args := u.Query()
			for k, v := range args {
				context.Set(key(k), v[0])
			}
		}
		// 是否需要控制参数
		err = checkArgs(context, target)
		if err != nil {
			context.Error(err)
			context.Status(400)
			context.Abort()
		}
		context.Next()
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
		// 字段名称
		value := fieldValueFromCtx(ctx, tp)
		if value == nil {
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

func fieldValueFromCtx(store mygin.ValueStore, tp reflect.StructField) interface{} {
	name := tp.Name
	// Tag名
	tag := tp.Tag.Get("json")
	// 寻找字段
	if temp := store.Value(key(strings.ToLower(name))); temp != nil {
		return temp
	} else if temp := store.Value(key(tag)); temp != nil {
		return temp
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
	if err != nil || n == 0 {
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

func checkArgs(ctx *mygin.Context, target interface{}) error {
	t := reflect.TypeOf(target).Elem()
	count := t.NumField()
	for i := 0; i < count; i++ {
		f := t.Field(i)
		tags := strings.Split(f.Tag.Get("args"), ";")
		value := fieldValueFromCtx(ctx, f)
		for _, tag := range tags {
			// 寻找参数
			err := getArgsHandleFunc(tag)(value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type argsHandleFunc func(value interface{}) error

func getArgsHandleFunc(tag string) argsHandleFunc {
	switch tag {
	case "required":
		return argsRequired
	default:
		// TODO：添加提示
		return argsDoNothing
	}
}

func argsRequired(value interface{}) error {
	if value == nil {
		return errors.New("参数缺失")
	}
	return nil
}

func argsDoNothing(value interface{}) error {
	return nil
}
