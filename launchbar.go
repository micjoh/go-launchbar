// package launchbar is a package to quickly write LaunchBar v6 actions like a pro
//
// For example check :
//   https://github.com/nbjahan/launchbar-pinboard
//   https://github.com/nbjahan/launchbar-spotlight
package launchbar

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/DHowett/go-plist"
	"github.com/codegangsta/inject"
)

type infoPlist map[string]interface{}

// Action represents a LaunchBar action
type Action struct {
	inject.Injector // Used for dependency injection
	Config          *Config
	Cache           Cache
	Input           *Input
	Logger          *log.Logger
	name            string
	views           map[string]*View
	items           []*Item
	context         *Context
	funcs           *FuncMap
	info            infoPlist
}

// NewAction creates an empty action, ready to populate with views
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

// Init parses the input
func (a *Action) Init(m ...FuncMap) *Action {
	a.funcs = &FuncMap{}
	if m != nil {
		*a.funcs = m[0]
	}

	in := NewInput(a, os.Args[1:])
	a.Input = in
	a.context.Input = in

	return a
}

// Run returns the compiled output of views. You must call Init first
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
						s := ""
						switch res := vals[0].Interface().(type) {

						case Items:
							s = res.Compile()
						case string:
							s = res
						case *View:
							s = res.Compile()
						case *Items:
							s = res.Compile()
						}
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
							var s string
							if out, ok := vals[0].Interface().(Items); ok {
								s = out.Compile()
							} else {
								s = vals[0].Interface().(*Items).Compile()
							}
							return s
						}
					}
					return ""
				}
			}
		}
	}

	// TODO: if a.GetView(view) == nil inform the user
	view := a.Config.GetString("view")
	if view == "" {
		view = "main"
	}
	w := a.GetView("*")
	out := a.GetView(view).Join(w).Compile()
	return out
}

// ShowView reruns the LaunchBar with the specified view.
//
// Use this when your LiveFeedback is enabled and you want to show another view
func (a *Action) ShowView(v string) {
	a.Config.Set("view", v)
	exec.Command("osascript", "-e", fmt.Sprintf(`tell application "LaunchBar"
       remain active
       perform action "%s"
       end tell`, a.name)).Start()
}

// NewView created a new view ready to populate with Items
func (a *Action) NewView(name string) *View {
	v := &View{a, name, make(Items, 0)}
	a.views[name] = v
	return v
}

// GetView returns a View if the view is not defined returns nil
func (a *Action) GetView(v string) *View {
	view, ok := a.views[v]
	if ok {
		return view
	}
	return nil
}

// GetItem return an Item with its ID. Returns nil if not found.
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

// Varsion returns Action version specified by CFBundleVersion key in Info.plist
func (a *Action) Version() Version { return Version(a.info["CFBundleVersion"].(string)) }

// LaunchBar provided variabled

// ActionPath returns the absolute path to the .lbaction bundle.
func (a *Action) ActionPath() string { return os.Getenv("LB_ACTION_PATH") }

// CachePath returns the absolute path to the action’s cache directory:
//  ~/Library/Caches/at.obdev.LaunchBar/Actions/Action Bundle Identifier/
//
// The action’s cache directory can be used to store files that can be recreated
// by the action itself, e.g. by downloading a file from a server again.
//
// Currently, this directory’s contents will never be touched by LaunchBar,
// but it may be periodically cleared in a future release.
// When the action is run, this directory is guaranteed to exist.
func (a *Action) CachePath() string { return os.Getenv("LB_CACHE_PATH") }

// Supportpath returns the The absolute path to the action’s support directory:
//  ~/Library/Application Support/LaunchBar/Action Support/Action Bundle Identifier/
//
// The action support directory can be used to persist user data between runs of
// the action, like preferences. When the action is run, this directory is
// guaranteed to exist.
func (a *Action) SupportPath() string { return os.Getenv("LB_SUPPORT_PATH") }

// IsDebug returns the value corresponds to LBDebugLogEnabled in the action’s Info.plist.
func (a *Action) IsDebug() bool { return os.Getenv("LB_DEBUG_LOG_ENABLED") == "true" }

// Launchbarpath returns the path to the LaunchBar.app bundle.
func (a *Action) LaunchBarPath() string { return os.Getenv("LB_LAUNCHBAR_PATH") }

// ScriptType returns the type of the script, as defined by the action’s Info.plist.
//
// This is either “default”, “suggestions” or “actionURL”.
//
// See http://www.obdev.at/resources/launchbar/developer-documentation/action-programming-guide.html#script-types for more information.
func (a *Action) ScriptType() string { return os.Getenv("LB_SCRIPT_TYPE") }

// IsCommandKey returns true if the Command key was down while running the action.
func (a *Action) IsCommandKey() bool { return os.Getenv("LB_OPTION_COMMAND_KEY") == "1" }

// IsOptionKey returns true if the Alternate (Option) key was down while running the action.
func (a *Action) IsOptionKey() bool { return os.Getenv("LB_OPTION_ALTERNATE_KEY") == "1" }

// IsShiftKey returns true if the Shift key was down while running the action.
func (a *Action) IsShiftKey() bool { return os.Getenv("LB_OPTION_SHIFT_KEY") == "1" }

// IsControlKey returns true if the Control key was down while running the action.
func (a *Action) IsControlKey() bool { return os.Getenv("LB_OPTION_CONTROL_KEY") == "1" }

// IsBackground returns true if the action is running in background.
func (a *Action) IsBackground() bool { return os.Getenv("LB_OPTION_RUN_IN_BACKGROUND") == "1" }
