package launchbar

import "encoding/json"

type Func interface{}
type FuncMap map[string]Func

var AlwasMatch = func() bool { return true }
var NeverMatch = func() bool { return false }
var MatchIfTrueFunc = func(b bool) func() bool { return func() bool { return b } }
var MatchIfFalseFunc = func(b bool) func() bool { return func() bool { return !b } }
var ShowViewFunc = func(v string) func(*Context) { return func(c *Context) { c.Action.ShowView(v) } }

// Item represents the LaunchBar item
type Item struct {
	View     *View
	item     *item
	match    Func
	run      Func
	render   Func
	children []item
}

func NewItem(title string) *Item {
	return &Item{
		item: &item{
			Title: title,
			Data:  make(map[string]interface{}),
		},
	}
}

func newItem(item *item) *Item {
	return &Item{item: item}
}

type item struct {
	// Standard fields
	Title                  string  `json:"title,omitempty"`
	Subtitle               string  `json:"subtitle,omitempty"`
	URL                    string  `json:"url,omitempty"`
	Path                   string  `json:"path,omitempty"`
	Icon                   string  `json:"icon,omitempty"`
	QuickLookURL           string  `json:"quickLookURL,omitempty"`
	Action                 string  `json:"action,omitempty"`
	ActionArgument         string  `json:"actionArgument,omitempty"`
	ActionReturnsItems     bool    `json:"actionReturnsItems,omitempty"`
	ActionRunsInBackground bool    `json:"actionRunsInBackground,omitempty"`
	ActionBundleIdentifier string  `json:"actionBundleIdentifier,omitempty"`
	Children               []*item `json:"children,omitempty"`

	// Custom fields
	ID       int                    `json:"x-id,omitempty"`
	Order    int                    `json:"x-order,omitempty"`
	FuncName string                 `json:"x-func,omitempty"`
	FuncArg  string                 `json:"x-funcarg,omitempty"`
	Arg      string                 `json:"x-arg,omitempty"`
	Data     map[string]interface{} `json:"x-data,omitempty"`
}

func (i *Item) SetTitle(title string) *Item              { i.item.Title = title; return i }
func (i *Item) SetSubtitle(subtitle string) *Item        { i.item.Subtitle = subtitle; return i }
func (i *Item) SetURL(url string) *Item                  { i.item.URL = url; return i }
func (i *Item) SetPath(path string) *Item                { i.item.Path = path; return i }
func (i *Item) SetIcon(icon string) *Item                { i.item.Icon = icon; return i }
func (i *Item) SetQuickLookURL(qlurl string) *Item       { i.item.QuickLookURL = qlurl; return i }
func (i *Item) SetAction(action string) *Item            { i.item.Action = action; return i }
func (i *Item) SetActionArgument(arg string) *Item       { i.item.ActionArgument = arg; return i }
func (i *Item) SetActionBundleIdentifier(s string) *Item { i.item.ActionBundleIdentifier = s; return i }
func (i *Item) SetActionRunsInBackground(b bool) *Item   { i.item.ActionRunsInBackground = b; return i }
func (i *Item) SetActionReturnsItems(b bool) *Item       { i.item.ActionReturnsItems = b; return i }
func (i *Item) SetChildren(items *Items) *Item           { i.item.Children = items.getItems(); return i }
func (i *Item) SetMatch(fn Func) *Item                   { i.match = fn; return i }
func (i *Item) SetRun(fn Func) *Item                     { i.run = fn; return i }
func (i *Item) SetRender(fn Func) *Item                  { i.render = fn; return i }
func (i *Item) SetOrder(n int) *Item                     { i.item.Order = n; return i }
func (i *Item) Item() *item                              { return i.item }
func (i *Item) Done() *View                              { return i.View }

func (i *Item) Run(f string, args ...interface{}) *Item {
	i.item.FuncName = f
	var ok bool
	var s string
	if len(args) == 1 {
		if s, ok = args[0].(string); ok {
			i.item.FuncArg = s
		}
	}
	if len(args) > 1 || !ok {
		b, err := json.Marshal(args)
		if err == nil {
			i.item.FuncArg = string(b)
		}
	}
	return i
}
