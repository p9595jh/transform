package transform

import (
	"fmt"
	"reflect"
	"strings"
)

type transformer struct {
	transformers map[string]f
	tag          string
}

func New(transformers ...I) Transformer {
	m := make(map[string]f, len(transformers)+len(defaults))
	for _, i := range defaults {
		m[i.Name] = i.F
	}
	for _, i := range transformers {
		m[i.Name] = i.F
	}
	return &transformer{m, _TAG}
}

func (t *transformer) RegisterTransformer(name string, f f) {
	t.transformers[name] = f
}

func (t *transformer) SetTag(tag string) Transformer {
	t.tag = tag
	return t
}

func (t *transformer) Tag() string {
	return t.tag
}

func (t *transformer) transform(value *reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		innerValue := reflect.Indirect(*value).FieldByName(field.Name)
		if field.Type.Kind() == reflect.Struct {
			t.transform(&innerValue)
		} else if tagString, ok := field.Tag.Lookup(t.tag); ok {
			if innerValue.CanSet() {
				for _, eachTag := range strings.Split(tagString, _SEPERATOR) {
					kv := strings.Split(eachTag, ":")
					var (
						tagName      = kv[0]
						tagParameter string
					)
					if len(kv) > 1 {
						tagParameter = kv[1]
					}
					if work, ok := t.transformers[tagName]; ok {
						worked := work(&shell{innerValue.Interface()}, tagParameter)
						innerValue.Set(reflect.ValueOf(worked.v))
					}
				}
			}
		}
	}
}

func (t *transformer) Transform(a any) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	elem := reflect.ValueOf(a).Elem()
	t.transform(&elem)
	return
}

// Field names of src and dst must be same.
func (t *transformer) mapping(srcValue, dstValue *reflect.Value) {
	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		srcInnerValue := reflect.Indirect(*srcValue).FieldByName(field.Name)
		dstInnerValue := reflect.Indirect(*dstValue).FieldByName(field.Name)
		if field.Type.Kind() == reflect.Struct {
			t.mapping(&srcInnerValue, &dstInnerValue)
		} else if tagString, ok := field.Tag.Lookup(t.tag); ok {
			if srcInnerValue.CanSet() && dstInnerValue.CanSet() {
				processing := &shell{srcInnerValue.Interface()}
				for _, eachTag := range strings.Split(tagString, _SEPERATOR) {
					kv := strings.Split(eachTag, ":")
					var (
						tagName      = kv[0]
						tagParameter string
					)
					if len(kv) > 1 {
						tagParameter = kv[1]
					}
					if work, ok := t.transformers[tagName]; ok {
						processing = &shell{work(processing, tagParameter).v}
					}
				}
				dstInnerValue.Set(reflect.ValueOf(processing.v))
			}
		}
	}
}

func (t *transformer) Mapping(src, dst any) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	srcElem := reflect.ValueOf(src).Elem()
	dstElem := reflect.ValueOf(dst).Elem()
	t.mapping(&srcElem, &dstElem)
	return
}
