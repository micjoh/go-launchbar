package launchbar

import "log"

type Context struct {
	Action *Action
	Config Config
	Cache  Cache
	Self   *Item
	Input  *Input
	Logger *log.Logger
}
