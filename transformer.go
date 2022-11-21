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
	t := &transformer{m, _TAG}
	for _, i := range transformers {
		t.add(i.Name, i.F)
	}
	return t
}

func (t *transformer) RegisterTransformer(name string, f f) {
	t.add(name, f)
}

func (t *transformer) add(name string, f f) {
	if strings.Contains(name, ",") {
		panic("cannot contain comma")
	}
	if strings.Contains(name, ":") {
		panic("cannot contain colon")
	}
	for _, s := range unavailables {
		if name == s {
			panic(name + " is not available (reserved word)")
		}
	}
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

func (t *transformer) mapping(paths []string, srcValue, dstValue *reflect.Value) {
mainLoop:
	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		eachPaths := append(paths, field.Name)
		srcInnerValue := reflect.Indirect(*srcValue).FieldByName(field.Name)
		if field.Type.Kind() == reflect.Struct {
			t.mapping(eachPaths, &srcInnerValue, dstValue)
		} else {
			var (
				finalPaths []string
				finalValue reflect.Value
			)
			if tagString, ok := field.Tag.Lookup(t.tag); ok {
				var (
					processing = &shell{srcInnerValue.Interface()}
					target     string
				)
				for _, eachTag := range strings.Split(tagString, _SEPERATOR) {
					var (
						kv           = strings.Split(eachTag, ":")
						tagName      = kv[0]
						tagParameter string
					)
					if len(kv) > 1 {
						tagParameter = kv[1]
						if tagName == _MAPPING {
							if tagParameter == _SKIP {
								continue mainLoop
							}
							target = tagParameter
						}
					}
					if work, ok := t.transformers[tagName]; ok {
						processing = &shell{work(processing, tagParameter).v}
					}
				}
				if target == "" {
					finalPaths = eachPaths
				} else {
					targetPaths := strings.Split(target, ".")
					if len(targetPaths) == 0 {
						panic("mapping target is empty: " + field.Name)
					}
					finalPaths = targetPaths
				}
				finalValue = reflect.ValueOf(processing.v)
			} else {
				finalPaths = eachPaths
				finalValue = srcInnerValue
			}
			dstTargetValue := reflect.Indirect(*dstValue).FieldByName(finalPaths[0])
			for i := 1; i < len(finalPaths); i++ {
				dstTargetValue = dstTargetValue.FieldByName(finalPaths[i])
			}
			if dstTargetValue.CanSet() {
				dstTargetValue.Set(finalValue)
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
	t.mapping([]string{}, &srcElem, &dstElem)
	return
}
