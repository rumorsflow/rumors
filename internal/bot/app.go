package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iagapie/rumors/internal/config"
	"github.com/iagapie/rumors/internal/notifications"
	"github.com/iagapie/rumors/pkg/logger"
	"github.com/rs/zerolog"
	"sync"
)

type Handler interface {
	Process(update tgbotapi.Update) error
}

type App struct {
	notification notifications.Notification
	log          *zerolog.Logger
	bot          *tgbotapi.BotAPI
	mu           sync.Mutex
}

func NewApp(debug bool, cfg config.TelegramConfig, log *zerolog.Logger) (*App, error) {
	_ = tgbotapi.SetLogger(logger.NewBotLogger(log))

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	return &App{
		notification: notifications.NewTgBotNotification(cfg.Owner, bot, log),
		log:          log,
		bot:          bot,
	}, nil
}

func (a *App) Bot() *tgbotapi.BotAPI {
	return a.bot
}

func (a *App) Notification() notifications.Notification {
	return a.notification
}

func (a *App) Start(handler Handler) {
	a.log.Info().Msg("Start telegram bot")
	a.start(handler)
}

func (a *App) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.log.Info().Msg("Stop telegram bot")
	a.bot.StopReceivingUpdates()
}

func (a *App) start(handler Handler) {
	a.mu.Lock()
	bot := a.bot
	a.mu.Unlock()

	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 30

	updates := bot.GetUpdatesChan(uc)

	for update := range updates {
		if err := handler.Process(update); err != nil {
			a.notification.Err(nil, err)
		}
	}
}
