package mgo

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	methodResNum = 2

	OptIgnore    = "-"
	OptOmitempty = "omitempty"
	OptIn        = "in"
	OptVague     = "vague"
)
const (
	flagIgnore = 1 << iota
	flagOmiEmpty
	flagIn
	flagVague
)

// StructToBson convert a golang struct to a bson
func StructToBson(s interface{}, tag string, methodName string) (bsonD bson.D, err error) {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, fmt.Errorf("%s is a nil pointer", v.Kind().String())
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 接受类型必须为 struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("s is not a struct but %s", v.Kind().String())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		// ignore unexported field
		if fieldType.PkgPath != "" {
			continue
		}

		// 读取标签
		tagVal, flag := readTag(fieldType, tag)

		if flag&flagIgnore != 0 {
			continue
		}
		fieldValue := v.Field(i)
		if flag&flagOmiEmpty != 0 && fieldValue.IsZero() {
			continue
		}

		// 忽略结构体中的空字段
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// 获取类型
		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			if methodName != "" {
				_, ok := fieldValue.Type().MethodByName(methodName)
				if ok {
					key, value, err := callFunc(fieldValue, methodName)
					if err != nil {
						return nil, err
					}
					bsonD = append(bsonD, bson.E{Key: key, Value: value})
					continue
				}
			}
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue})
		case reflect.Struct:
			if methodName != "" {
				_, ok := fieldValue.Type().MethodByName(methodName)
				if ok {
					key, value, err := callFunc(fieldValue, methodName)
					if err != nil {
						return nil, err
					}
					bsonD = append(bsonD, bson.E{Key: key, Value: value})
					continue
				}
			}
			// recursive
			deepRes, deepErr := StructToBson(fieldValue.Interface(), tag, methodName)
			if deepErr != nil {
				return nil, deepErr
			}
			if flag&flagIn != 0 {
				bsonD = append(bsonD, deepRes...)
			} else {
				bsonD = append(bsonD, bson.E{Key: tagVal, Value: deepRes})
			}
		case reflect.Map:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue})
		case reflect.Chan:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue})
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Int()})
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Uint()})
		case reflect.Float32, reflect.Float64:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Float()})
		case reflect.String:
			if flag&flagVague != 0 {
				regexString := CheckRegex(fieldValue.String())
				bsonD = append(bsonD, bson.E{Key: tagVal, Value: bson.D{{Key: "$regex", Value: regexString}}})
			} else {
				bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.String()})
			}
		case reflect.Bool:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Bool()})
		case reflect.Complex64, reflect.Complex128:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Complex()})
		case reflect.Interface:
			bsonD = append(bsonD, bson.E{Key: tagVal, Value: fieldValue.Interface()})
		default:
		}
	}

	return
}

// readTag read tag from struct field
func readTag(f reflect.StructField, tag string) (string, int) {
	lookup, ok := f.Tag.Lookup(tag)

	fieldTag := ""
	flag := 0

	// no tag, use field name
	if !ok {
		return f.Name, flag
	}
	opts := strings.Split(lookup, ",")

	fieldTag = opts[0]
	if len(opts) == 1 && opts[0] == OptIgnore {
		flag |= flagIgnore
	}
	for i := 1; i < len(opts); i++ {
		switch opts[i] {
		case OptOmitempty:
			flag |= flagOmiEmpty
		case OptIn:
			flag |= flagIn
		case OptVague:
			flag |= flagVague
		}
	}
	return fieldTag, flag
}

// call function
func callFunc(fv reflect.Value, methodName string) (string, interface{}, error) {
	methodRes := fv.MethodByName(methodName).Call([]reflect.Value{})
	if len(methodRes) != methodResNum {
		return "", nil, fmt.Errorf("wrong method %s, should have 2 output: (string,interface{})", methodName)
	}
	if methodRes[0].Kind() != reflect.String {
		return "", nil, fmt.Errorf("wrong method %s, first output should be string", methodName)
	}
	key := methodRes[0].String()
	value := methodRes[1]
	return key, value.Interface(), nil
}

// CheckRegex 处理模糊查询时 特殊字符串问题
func CheckRegex(str string) string {
	//正则匹配出现的特殊字符串
	fbsArr := []string{"\\", "$", "(", ")", "*", "+", ".", "[", "]", "?", "^", "{", "}", "|"}
	for _, ch := range fbsArr {
		if StrContainers := strings.Contains(str, ch); StrContainers {
			str = strings.Replace(str, ch, "\\"+ch, -1)
		}
	}
	return str
}
