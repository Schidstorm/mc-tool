package cli

import (
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
)

func Configure(command *cobra.Command, cfg interface{}) {
	flags := command.PersistentFlags()
	tp := reflect.Indirect(reflect.ValueOf(cfg)).Type()

	for i := 0; i < tp.NumField(); i++ {
		fieldType := tp.Field(i)
		tag, ok := fieldType.Tag.Lookup("cli")
		if ok {
			switch fieldType.Type.Kind() {
			case reflect.Bool:
				defaultValue, _ := strconv.ParseBool(fieldType.Tag.Get("default"))
				usage := fieldType.Tag.Get("usage")
				flags.Bool(tag, defaultValue, usage)
			case reflect.Int:
				var defaultValue int64
				d, ok := fieldType.Tag.Lookup("default")
				if ok {
					defaultValue, _ = strconv.ParseInt(d, 10, int(reflect.TypeOf(0).Size())*8)
				} else {
					defaultValue = 0
				}
				usage := fieldType.Tag.Get("usage")
				flags.Int(tag, int(defaultValue), usage)
			case reflect.String:
				defaultValue := fieldType.Tag.Get("default")
				usage := fieldType.Tag.Get("usage")
				flags.String(tag, defaultValue, usage)
			default:
				panic("not implemented cli config type")
			}
		}
	}
}

func GetValues(command *cobra.Command, cfg interface{}) error {
	flags := command.PersistentFlags()
	val := reflect.ValueOf(cfg).Elem()
	tp := val.Type()

	for i := 0; i < tp.NumField(); i++ {
		fieldType := tp.Field(i)
		tag, ok := fieldType.Tag.Lookup("cli")
		if ok {
			switch fieldType.Type.Kind() {
			case reflect.Bool:
				if v, err := flags.GetBool(tag); err != nil {
					return err
				} else {
					val.Field(i).SetBool(v)
				}
			case reflect.Int:
				if v, err := flags.GetInt(tag); err != nil {
					return err
				} else {
					val.Field(i).SetInt(int64(v))
				}
			case reflect.String:
				if v, err := flags.GetString(tag); err != nil {
					return err
				} else {
					valField := val.FieldByName(fieldType.Name)
					if valField.CanSet() {
						valField.SetString(v)
					}
				}
			default:
				panic("not implemented cli config type")
			}
		}
	}

	return nil
}
