package launchbar

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/DHowett/go-plist"
	"github.com/codegangsta/inject"
)

type infoPlist map[string]interface{}

type Action struct {
	inject.Injector
	name    string
	Config  Config
	Cache   Cache
	views   map[string]*View
	items   []*Item
	Input   *Input
	Logger  *log.Logger
	context *Context
	funcs   *FuncMap
	info    infoPlist
}

func NewAction(name string, config ConfigValues) *Action {
	a := &Action{
		Injector: inject.New(),
		name:     name,
		views:    make(map[string]*View),
		items:    make([]*Item, 0),
	}
	a.Config = NewConfigDefaults(a.SupportPath(), config)
	a.Cache = NewCache(a.CachePath())
	fd, err := os.OpenFile(path.Join(a.SupportPath(), "error.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		fd = os.Stderr
	}
	a.Logger = log.New(fd, "", 0)
	c := &Context{
		Action: a,
		Config: a.Config,
		Cache:  a.Cache,
		Logger: a.Logger,
	}
	a.context = c
	a.Map(c)

	data, err := ioutil.ReadFile(path.Join(a.ActionPath(), "Contents", "Info.plist"))
	if err != nil {
		a.Logger.Println(err)
		panic(err)
	}
	_, err = plist.Unmarshal(data, &a.info)
	if err != nil {
		a.Logger.Println(err)
		panic(err)
	}
	return a
}

func (a *Action) Init(m ...FuncMap) *Action {
	a.funcs = &FuncMap{}
	if m != nil {
		*a.funcs = m[0]
	}

	in := NewInput(a, strings.Join(os.Args[1:], " "))
	a.Input = in
	a.context.Input = in

	return a
}

func (a *Action) Run() string {
	in := a.Input
	if in.IsObject() {
		if in.hasFunc {
			// I'm not sure!
			a.context.Self = in.Item
			if fn, ok := (*a.funcs)[in.Item.item.FuncName]; ok {
				vals, err := a.Invoke(fn)
				if err != nil {
					a.Logger.Fatalln(err)
				}
				if len(vals) > 0 {
					if vals[0].Interface() != nil {
						out := vals[0].Interface().(Items)
						s := out.Compile()
						return s
					}
				}
				return ""
			}
		} else {
			if item := a.GetItem(in.Item.item.ID); item != nil {
				a.context.Self = item
				if item.run != nil {
					vals, err := a.Invoke(item.run)
					if err != nil {
						a.Logger.Fatalln(err)
					}
					if len(vals) > 0 {
						if vals[0].Interface() != nil {
							out := vals[0].Interface().(Items)
							s := out.Compile()
							return s
						}
					}
					return ""
				}
			}
		}
	}
	view := a.Config.GetString("view")
	if view == "" {
		view = "main"
	}
	w := a.GetView("*")
	out := a.GetView(view).Join(w).Compile()
	return out
}

func (a *Action) ShowView(v string) {
	a.Config.Set("view", v)
	exec.Command("osascript", "-e", fmt.Sprintf(`tell application "LaunchBar"
       remain active
       perform action "%s"
       end tell`, a.name)).Start()
}

func (a *Action) NewView(name string) *View {
	v := &View{a, name, make(Items, 0)}
	a.views[name] = v
	return v
}

func (a *Action) GetView(v string) *View {
	view, ok := a.views[v]
	if ok {
		return view
	}
	return nil
}

func (a *Action) GetItem(id int) *Item {
	if id < 1 {
		return nil
	}
	if id > len(a.items) {
		return nil
	}
	return a.items[id-1]
}

// Info.plist variables
func (a *Action) Version() Version { return Version(a.info["CFBundleVersion"].(string)) }

// LauncBar provided variabled
func (a *Action) ActionPath() string    { return os.Getenv("LB_ACTION_PATH") }
func (a *Action) CachePath() string     { return os.Getenv("LB_CACHE_PATH") }
func (a *Action) SupportPath() string   { return os.Getenv("LB_SUPPORT_PATH") }
func (a *Action) IsDebug() bool         { return os.Getenv("LB_DEBUG_LOG_ENABLED") == "true" }
func (a *Action) LaunchBarPath() string { return os.Getenv("LB_LAUNCHBAR_PATH") }
func (a *Action) ScriptType() string    { return os.Getenv("LB_SCRIPT_TYPE") }
func (a *Action) IsCommandKey() bool    { return os.Getenv("LB_OPTION_COMMAND_KEY") == "1" }
func (a *Action) IsOptionKey() bool     { return os.Getenv("LB_OPTION_ALTERNATE_KEY") == "1" }
func (a *Action) IsShiftKey() bool      { return os.Getenv("LB_OPTION_SHIFT_KEY") == "1" }
func (a *Action) IsControlKey() bool    { return os.Getenv("LB_OPTION_CONTROL_KEY") == "1" }
func (a *Action) ISBackground() bool    { return os.Getenv("LB_OPTION_RUN_IN_BACKGROUND") == "1" }
