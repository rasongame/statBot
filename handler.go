package main

import "statBot/utils"

func AddHandler(command string, handler utils.HandlerFunc, filter utils.FilterFunc) (utils.Handler, bool) {
	if h, ok := utils.Handlers[command]; ok {
		if !ok {
			return h, false
		}
	}
	h := utils.Handler{handler, filter}
	utils.Handlers[command] = h
    utils.Handler
    return h, true
}
