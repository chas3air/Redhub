package app

import "log/slog"

type App struct {
}

func New(log *slog.Logger, port int, storage any) *App {
	
	return &App{}
}
