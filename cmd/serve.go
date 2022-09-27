package cmd

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/background"
	backHandlers "github.com/iagapie/rumors/internal/background/handlers"
	"github.com/iagapie/rumors/internal/bot"
	botHandlers "github.com/iagapie/rumors/internal/bot/handlers"
	"github.com/iagapie/rumors/internal/config"
	"github.com/iagapie/rumors/internal/storage"
	"github.com/iagapie/rumors/pkg/logger"
	"github.com/iagapie/rumors/pkg/mongodb"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	golog "log"
	"os/signal"
	"syscall"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.LocalFlags())
	},
	RunE: serve,
}

func init() {
	flagSet := serveCmd.PersistentFlags()

	flagSet.String("server.network", "tcp", "server network, ex: tcp, tcp4, tcp6, unix, unixpacket")
	flagSet.String("server.address", ":8080", "server address")

	flagSet.String("mongodb.uri", "", "mongo db uri")

	flagSet.Int64("telegram.owner", 0, "telegram (BOT) owner id")
	flagSet.String("telegram.token", "", "telegram bot token")

	flagSet.String("async.redis.network", "tcp", "redis network, ex: tcp, unix")
	flagSet.String("async.redis.address", ":6379", "redis address")
	flagSet.String("async.redis.username", "", "redis username")
	flagSet.String("async.redis.password", "", "redis password")
	flagSet.Int("asynq.redis.db", 0, "by default, redis offers 16 databases (0..15)")
	flagSet.Int("asynq.server.concurrency", 0, "how many concurrent workers to use, zero or negative for number of CPUs")
	flagSet.String("asynq.scheduler.cron", "@every 5m", "can use cron spec string or can use \"@every <duration>\" to specify the interval")
	flagSet.String("asynq.scheduler.name", "feeds:parser", "task type (handler route)")

	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, _ []string) error {
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	log := zeroLog(cmd, cfg)
	golog.SetOutput(log)

	botLog := log.With().Str("context", "bot").Logger()
	backLog := log.With().Str("context", "back").Logger()
	//httpLog := log.With().Str("context", "http").Logger()

	mongoDB, err := mongodb.GetDB(cmd.Context(), cfg.MongoDB.URI)
	if err != nil {
		return err
	}

	roomStorage, err := storage.NewRoomStorage(cmd.Context(), mongoDB)
	if err != nil {
		return err
	}

	feedStorage, err := storage.NewFeedStorage(cmd.Context(), mongoDB)
	if err != nil {
		return err
	}

	feedItemStorage, err := storage.NewFeedItemStorage(cmd.Context(), mongoDB)
	if err != nil {
		return err
	}

	backApp := background.NewApp(cfg.Asynq, &backLog)
	//httpApp := http.NewApp(cfg.Debug, cfg.Server, &httpLog)
	botApp, err := bot.NewApp(cfg.Debug, cfg.Telegram, &botLog)
	if err != nil {
		return err
	}

	roomsHandler := &backHandlers.RoomsHandler{
		Notification: botApp.Notification(),
		RoomStorage:  roomStorage,
		Client:       backApp.Client(),
		Log:          &backLog,
	}

	feedsHandler := &backHandlers.FeedsHandler{
		Notification: botApp.Notification(),
		FeedStorage:  feedStorage,
		Client:       backApp.Client(),
		Log:          &backLog,
	}

	roomsMux := asynq.NewServeMux()
	roomsMux.Handle("rooms:crud", roomsHandler)
	roomsMux.Handle("rooms:list", roomsHandler)
	roomsMux.Handle("rooms:view", roomsHandler)
	roomsMux.Handle("rooms:update", roomsHandler)
	roomsMux.Handle("rooms:add", &backHandlers.RoomsAddHandler{
		Notification: botApp.Notification(),
		RoomStorage:  roomStorage,
		Log:          &backLog,
	})
	roomsMux.Handle("rooms:left", &backHandlers.RoomsLeftHandler{
		Notification: botApp.Notification(),
		RoomStorage:  roomStorage,
		Log:          &backLog,
	})

	feedsMux := asynq.NewServeMux()
	feedsMux.Handle("feeds:crud", feedsHandler)
	feedsMux.Handle("feeds:list", feedsHandler)
	feedsMux.Handle("feeds:view", feedsHandler)
	feedsMux.Handle("feeds:update", feedsHandler)
	feedsMux.Handle("feeds:add", &backHandlers.FeedsAddHandler{
		Notification: botApp.Notification(),
		FeedStorage:  feedStorage,
		Log:          &backLog,
		Owner:        cfg.Telegram.Owner,
	})
	feedsMux.Handle(cfg.Asynq.Scheduler.TaskName, &backHandlers.FeedsParserHandler{
		Notification:    botApp.Notification(),
		FeedStorage:     feedStorage,
		FeedItemStorage: feedItemStorage,
		Client:          backApp.Client(),
		Log:             &backLog,
	})

	rumorsMux := asynq.NewServeMux()
	rumorsMux.Handle("rumors:list", &backHandlers.RumorsHandler{
		Notification:    botApp.Notification(),
		FeedStorage:     feedStorage,
		FeedItemStorage: feedItemStorage,
		Log:             &backLog,
	})

	backMux := asynq.NewServeMux()
	backMux.Handle("rooms:", roomsMux)
	backMux.Handle("feeds:", feedsMux)
	backMux.Handle("rumors:", rumorsMux)
	backMux.Handle("aggregated:broadcast", &backHandlers.BroadcastHandler{
		Notification: botApp.Notification(),
		FeedStorage:  feedStorage,
		RoomStorage:  roomStorage,
		Log:          &backLog,
	})

	botHandler := &botHandlers.UpdateHandler{
		Notification: botApp.Notification(),
		Client:       backApp.Client(),
		Log:          &botLog,
		Owner:        cfg.Telegram.Owner,
	}

	//api := httpApp.Echo().Group("/api")
	//{
	//	v1 := api.Group("/v1")
	//	{
	//		new(httpHandlers.UpdateHandler).Register(v1)
	//	}
	//}

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	defer cancel()

	go func(ctx context.Context, backApp *background.App) {
		<-ctx.Done()

		backApp.Shutdown()
	}(ctx, backApp)

	//go func(ctx context.Context, httpApp *http.App) {
	//	<-ctx.Done()
	//
	//	httpApp.Shutdown()
	//}(ctx, httpApp)

	go func(ctx context.Context, botApp *bot.App) {
		<-ctx.Done()

		botApp.Shutdown()
	}(ctx, botApp)

	if err = backApp.Start(backMux); err != nil {
		return err
	}

	go botApp.Start(botHandler)

	<-ctx.Done()

	//return httpApp.Start()
	return nil
}

func zeroLog(cmd *cobra.Command, cfg config.Config) *zerolog.Logger {
	factory := &logger.ZeroLog{
		Command:  name(cmd),
		Version:  cmd.Version,
		Debug:    cfg.Debug,
		LogLevel: cfg.Log.Level,
		Colored:  cfg.Log.Colored,
	}
	return factory.Log()
}

func name(cmd *cobra.Command) string {
	if cmd.Parent() == nil {
		return cmd.Name()
	}
	return fmt.Sprintf("%s %s", name(cmd.Parent()), cmd.Name())
}
