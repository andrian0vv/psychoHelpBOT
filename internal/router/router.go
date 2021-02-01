package router

import (
	"sync"
)

type router struct {
	MessageHandler  Handler
	mux             sync.Mutex
	CommandHandlers map[string]Handler
}

func New() Router {
	return &router{
		CommandHandlers: make(map[string]Handler),
	}
}

func (r *router) RegisterMessageHandler(handler Handler) {
	r.MessageHandler = handler
}

func (r *router) RegisterCommandHandler(command string, handler Handler) {
	r.CommandHandlers[command] = handler
}

func (r *router)  GetMessageHandler() Handler {
	if r.MessageHandler != nil {
		return r.MessageHandler
	}

	return &DefaultHandler{}
}

func (r *router) GetCommandHandler(command string) (Handler, bool) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if handler, ok := r.CommandHandlers[command]; ok {
		return handler, true
	}

	return nil, false
}
