package launchbar

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Input represents the object that LaunchBar passes to scripts
type Input struct {
	raw      string
	Item     *Item
	isString bool
	isNumber bool
	hasFunc  bool
	isFloat  bool
	isInt    bool
	isObject bool
	isPaths  bool
	hasData  bool
	paths    []string
	number   float64
}

func NewInput(a *Action, s string) *Input {
	item := item{}
	var in = &Input{
		raw: s,
	}

	if err := json.Unmarshal([]byte(in.Raw()), &item); err == nil {
		in.isObject = true
		if item.Data != nil && len(item.Data) > 0 {
			in.hasData = true
		}
		in.Item = a.GetItem(item.ID)
		if in.Item == nil {
			in.Item = newItem(&item)
		}
		in.Item.item.Arg = item.Arg

		in.Item.item.Order = item.Order
		in.Item.item.FuncName = item.FuncName
		in.Item.item.Data = item.Data
		if item.FuncName != "" {
			in.hasFunc = true
		}
	} else {
		//TODO: Input > check for paths
		in.isString = true
	}

	if f64, err := strconv.ParseFloat(in.String(), 64); err == nil {
		in.isNumber = true
		in.number = f64
		if fmt.Sprintf("%f", f64) == fmt.Sprintf("%f", float64(int64(f64))) {
			in.isInt = true
		} else {
			in.isFloat = true
		}
	}

	return in

}

func (in *Input) Int() int         { return int(in.number) }
func (in *Input) Float64() float64 { return in.number }
func (in *Input) Int64() int64     { return int64(in.number) }
func (in *Input) Raw() string      { return in.raw }

func (in *Input) String() string {
	if in.IsObject() {
		return in.Item.item.Arg
	}
	return in.Raw()
}

func (in *Input) FuncArg() string {
	if in.IsObject() {
		return in.Item.item.FuncArg
	}
	return ""
}

func (in *Input) FuncArgsString() []string {
	if !in.isObject {
		return nil
	}
	var out []string
	err := json.Unmarshal([]byte(in.Item.item.FuncArg), &out)
	if err != nil {
		return []string{in.Item.item.FuncArg}
	}
	return out
}

func (in *Input) FuncArgsMapString() map[int]string {
	out := make(map[int]string)
	if !in.isObject {
		return out
	}
	var args []interface{}
	err := json.Unmarshal([]byte(in.Item.item.FuncArg), &args)
	if err != nil {
		out[0] = in.Item.item.FuncArg
		return out
	}
	for i, arg := range args {
		out[i] = fmt.Sprintf("%v", arg)
	}
	return out
}

func (in *Input) Title() string {
	if in.IsObject() {
		return in.Item.item.Title
	}
	return ""
}

func (in *Input) Data(key string) interface{} {
	if in.hasData {
		if i, ok := in.Item.item.Data[key]; ok {
			return i
		}
	}
	return nil
}

//DataString returns a customdata[key] as string
func (in *Input) DataString(key string) string {
	if in.Item == nil {
		return ""
	}
	if s, ok := in.Item.item.Data[key].(string); ok {
		return s
	}
	return ""
}

//DataInt returns a customdata[key] as string
func (in *Input) DataInt(key string) int {
	if s, ok := in.Item.item.Data[key].(int); ok {
		return s
	}
	return 0
}

func (in *Input) IsString() bool { return in.isString }
func (in *Input) IsObject() bool { return in.isObject }
func (in *Input) IsPaths() bool  { return in.isPaths }
func (in *Input) IsNumber() bool { return in.isNumber }
func (in *Input) IsInt() bool    { return in.isInt }
func (in *Input) IsFloat() bool  { return in.isFloat }
func (in *Input) IsEmpty() bool  { return in.String() == "" }
