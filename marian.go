package utils

import (
	"github.com/fatih/structtag"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

func CleanStruct(specificTag, optionInTag string, structWithTag any) any {

	if reflect.TypeOf(structWithTag).Kind() != reflect.Ptr {
		log.Warning("value must be a pointer")
		panic("value must be a pointer")
	}

	mt := reflect.ValueOf(structWithTag)
	dv := mt.Elem() // follow pointer

	if dv.Kind() != reflect.Struct {
		log.Warning("value must be a pointer to a struct/interface")
		panic("value must be a pointer to a struct/interface")
	}

	ParseVisibility(specificTag, optionInTag, dv)

	return structWithTag

}

func ParseVisibility(specificTag, optionInTag string, dv reflect.Value) any {

	for i := 0; i < dv.NumField(); i++ {
		log.Debug("Field :" + strconv.Itoa(i) + "," + dv.Type().Field(i).Name + " type:" + dv.Type().Field(i).Type.Kind().String())

		// get field tag
		tag := dv.Type().Field(i).Tag

		// ... and start using structtag by parsing the tag
		tags, err := structtag.Parse(string(tag))
		if err != nil {
			log.Warning("No parse tags:", err)
			panic(err)
		}

		// get a single tag
		workingThisTag, err := tags.Get(specificTag)
		if err != nil {
			log.Error(err.Error() + " -> " + specificTag + " in " + dv.Type().Field(i).Name)
			continue
		}

		containOption := containOption(optionInTag, workingThisTag.Name, workingThisTag.Options)
		if !containOption {
			log.Debug("Empty value for field:", dv.Type().Field(i).Name)
			emptyValue := reflect.Zero(dv.Field(i).Type())
			switch dv.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dv.Field(i).SetInt(emptyValue.Int())
			case reflect.String:
				dv.Field(i).SetString(emptyValue.String())
			case reflect.Bool:
				dv.Field(i).SetBool(emptyValue.Bool())
			case reflect.Float32, reflect.Float64:
				dv.Field(i).SetFloat(emptyValue.Float())
			case reflect.Struct:
				dv.Field(i).Set(emptyValue)
			case reflect.Slice:
				dv.Field(i).Set(emptyValue)
			default:
				{
					ouput := "Unsupported : " + dv.Field(i).Type().String()
					log.Error(ouput)
					panic(ouput)
				}
			}
		} else {
			switch dv.Field(i).Kind() {
			case reflect.Struct:
				{
					ParseVisibility(specificTag, optionInTag, dv.Field(i))
				}
			case reflect.Array:
				{
					log.Warning("is array")
				}
			case reflect.Slice:
				{
					log.Debug("is slice")
					switch dv.Field(i).Interface().(type) {
					case []string, []int, []int32, []int64, []bool, []float64, []float32:
						log.Debug("Primitive type")
					default:
						{
							log.Warning("Maybe slice struct:", dv.Field(i).Type())
							tempSlice := dv.Field(i)
							for i := 0; i < tempSlice.Len(); i++ {
								ParseVisibility(specificTag, optionInTag, tempSlice.Index(i))
							}
						}
					}

				}
			case reflect.Float32, reflect.Float64, reflect.String, reflect.Int, reflect.Int32, reflect.Int64, reflect.Bool:
				{
					log.Debug("nothing to make")
				}
			default:
				{
					log.Error("Type not supported:", dv.Field(i).Kind())
				}
			}
		}
	}
	return dv
}

func containOption(optionInTag, name string, options []string) bool {

	if optionInTag == name {
		return true
	}

	for i := range options {
		if optionInTag == options[i] {
			return true
		}
	}
	return false
}
