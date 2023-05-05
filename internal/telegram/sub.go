package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/conv"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"golang.org/x/exp/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	OpUnmarshalMessage  = "telegram sub: unmarshal message ->"
	OpPrepareMessage    = "telegram sub: prepare message ->"
	OpUnmarshalArticles = "telegram sub: unmarshal articles ->"
	OpFindSite          = "telegram sub: find site ->"
	OpFindChats         = "telegram sub: find chats ->"
	OpProcessor         = "telegram sub: processor ->"

	chatQuery = "field.0.0=broadcast&value.0.0=%s&field.1.0=blocked&value.1.0=false&field.2.0=deleted&value.2.0=false"
)

type Subscriber struct {
	pool     *pool
	bot      *Bot
	sub      *pubsub.Subscriber
	logger   *slog.Logger
	siteRepo repository.ReadRepository[*entity.Site]
	chatRepo repository.ReadRepository[*entity.Chat]
}

func NewSubscriber(
	bot *Bot,
	sub *pubsub.Subscriber,
	siteRepo repository.ReadRepository[*entity.Site],
	chatRepo repository.ReadRepository[*entity.Chat],
) *Subscriber {
	log := logger.WithGroup("telegram").WithGroup("subscriber")

	s := &Subscriber{
		bot:      bot,
		sub:      sub,
		logger:   log,
		siteRepo: siteRepo,
		chatRepo: chatRepo,
	}
	s.pool = newPool(s.processor, log.WithGroup("pool"))

	return s
}

func (s *Subscriber) Run(ctx context.Context) error {
	s.pool.Run(ctx)

	telegramSub := s.sub.Telegram(ctx)
	articlesSub := s.sub.Articles(ctx)

	telegramCh := telegramSub.Channel()
	articlesCh := articlesSub.Channel()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-telegramCh:
			if data == nil {
				continue
			}

			var message Message
			if err := json.Unmarshal(conv.StringToBytes(data.Payload), &message); err != nil {
				err = fmt.Errorf("%s error: %w", OpUnmarshalMessage, err)
				s.logger.Error("error due to unmarshal message", "err", err, "channel", data.Channel, "payload", data.Payload)
				continue
			}

			s.logger.Debug("message received", "channel", data.Channel, "message", message)

			s.send(message, data.Channel)

		case data := <-articlesCh:
			if data == nil {
				continue
			}

			var articles []pubsub.Article
			if err := json.Unmarshal(conv.StringToBytes(data.Payload), &articles); err != nil {
				err = fmt.Errorf("%s error: %w", OpUnmarshalArticles, err)
				s.logger.Error("error due to unmarshal articles", "err", err, "channel", data.Channel, "payload", data.Payload)
				continue
			}

			s.logger.Debug("articles received", "channel", data.Channel, "articles", articles)

			for i := len(articles) - 1; i >= 0; i-- {
				article := articles[i]

				site, err := s.siteRepo.FindByID(ctx, article.SiteID)
				if err != nil {
					err = fmt.Errorf("%s error: %w", OpFindSite, err)
					s.logger.Error("error due to find site", "err", err, "channel", data.Channel, "article", article.ID, "site", article.SiteID)
					continue
				}

				if !*site.Enabled {
					err = fmt.Errorf("%s site not found", OpFindSite)
					s.logger.Debug("error due to find site", "err", err, "channel", data.Channel, "article", article.ID, "site", article.SiteID)
					continue
				}

				chats, err := s.chatRepo.Find(ctx, db.BuildCriteria(fmt.Sprintf(chatQuery, article.SiteID)))
				if err != nil {
					err = fmt.Errorf("%s error: %w", OpFindChats, err)
					s.logger.Error("error due to find chats", "err", err, "channel", data.Channel, "article", article.ID, "site", article.SiteID)
					continue
				}

				if len(chats) == 0 {
					err = fmt.Errorf("%s chats not found", OpFindChats)
					s.logger.Debug("error due to find chats", "err", err, "channel", data.Channel, "article", article.ID, "site", article.SiteID)
					continue
				}

				for _, chat := range chats {
					s.send(Message{
						ChatID:   chat.TelegramID,
						ImageURL: article.Image,
						View:     ViewArticle,
						Data:     article,
						Delay:    true,
					}, data.Channel)
				}
			}
		}
	}
}

func (s *Subscriber) send(message Message, channel string) {
	chunks, err := message.chattable(s.bot)
	if err != nil {
		err = fmt.Errorf("%s error: %w", OpPrepareMessage, err)
		s.logger.Error("error due prepare message before send", "err", err, "channel", channel, "message", message)
		return
	}

	s.pool.Add(message.ChatID, newEntry(chunks, message.Delay))
}

func (s *Subscriber) processor(message tgbotapi.Chattable) error {
	if _, err := s.bot.Request(message); err != nil {
		var res []byte
		if e, ok := err.(*tgbotapi.Error); ok {
			res, _ = json.Marshal(e)
		}
		s.logger.Error("bot request error", "err", fmt.Errorf("%s error: %w", OpProcessor, err), "message", message, "response", res)

		return err
	}
	return nil
}

type pool struct {
	t         *time.Ticker
	max       int
	logger    *slog.Logger
	waiting   *Dict[[]*entry]
	workers   *Dict[*worker]
	processor func(tgbotapi.Chattable) error
}

func newPool(processor func(tgbotapi.Chattable) error, logger *slog.Logger) *pool {
	return &pool{
		max:       25,
		logger:    logger,
		waiting:   NewDict[[]*entry](),
		workers:   NewDict[*worker](),
		processor: processor,
	}
}

func (p *pool) Add(chatID int64, e *entry) {
	name := fmt.Sprintf("worker%d", chatID)

	if p.waiting.Has(name) {
		p.waiting.Set(name, append(p.waiting.Get(name), e))
	} else {
		p.waiting.Set(name, []*entry{e})
	}
}

func (p *pool) Run(ctx context.Context) {
	go p.run(ctx)
}

func (p *pool) run(ctx context.Context) {
	p.t = time.NewTicker(time.Second)

	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		p.t.Stop()
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.t.C:
			p.process(ctx)
		}
	}
}

func (p *pool) process(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	p.workers.Loop(func(key string, item *worker) bool {
		if item.Closed() {
			p.logger.Debug("remove worker", "worker", key)
			p.workers.Del(key)
		}
		return true
	})

	need := p.max - p.workers.Len()

	p.waiting.Loop(func(key string, item []*entry) bool {
		var w *worker

		if w = p.workers.Del(key); w != nil {
			if w.Closed() {
				need++
			} else {
				p.waiting.Del(key)
				p.logger.Debug("old worker found", "worker", key)
			}
		}

		if w == nil && need > 0 {
			need--

			w = newWorker(key, p.processor)
			p.logger.Debug("new worker", "worker", key)

			w.Run(ctx)
			p.logger.Debug("worker run", "worker", key)
		}

		if w != nil {
			p.workers.Set(key, w)

			for _, e := range item {
				w.Put(e)
			}
		}

		return true
	})
}

type entry struct {
	delay bool
	data  *List[*chattable]
}

type chattable struct {
	tgbotapi.Chattable
	max   uint // Max retry
	retry uint
}

func newEntry(data []tgbotapi.Chattable, delay bool) *entry {
	if len(data) == 0 {
		panic("newEntry empty data")
	}

	e := &entry{delay: delay, data: NewList[*chattable](len(data))}
	for _, item := range data {
		e.data.Add(&chattable{Chattable: item, max: 3, retry: 1})
	}

	return e
}

func (e *entry) process(processor func(tgbotapi.Chattable) error) int {
	for !e.data.IsEmpty() {
		item := e.data.Shift()

		if err := processor(item.Chattable); err != nil {
			if tge, ok := err.(*tgbotapi.Error); ok && tge.RetryAfter > 0 && item.retry < item.max {
				item.retry++

				e.data.Unshift(item)

				return tge.RetryAfter
			}

			return 0
		}
	}

	return 0
}

type worker struct {
	t         *time.Ticker
	d         time.Duration
	mu        sync.Mutex
	name      string
	entries   *List[*entry]
	running   atomic.Bool
	processor func(tgbotapi.Chattable) error
}

func newWorker(name string, processor func(tgbotapi.Chattable) error) *worker {
	w := &worker{
		d:         3 * time.Second,
		name:      name,
		entries:   NewList[*entry](0),
		processor: processor,
	}
	w.t = time.NewTicker(w.d)
	return w
}

func (w *worker) String() string {
	return w.name
}

func (w *worker) Put(e *entry) {
	if w.Closed() {
		return
	}

	if e.delay {
		w.entries.Add(e)
	} else {
		w.process(e)
	}
}

func (w *worker) Close() {
	w.running.Store(false)
}

func (w *worker) Closed() bool {
	return !w.running.Load()
}

func (w *worker) Run(ctx context.Context) {
	go w.run(ctx)
}

func (w *worker) run(ctx context.Context) {
	w.running.Store(true)

	defer func() {
		w.t.Stop()
		w.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.t.C:
			if w.entries.IsEmpty() {
				return
			}

			e := w.entries.Shift()

			if !w.process(e) {
				return
			}
		}
	}
}

func (w *worker) process(e *entry) bool {
	w.mu.Lock()
	w.t.Stop()
	w.mu.Unlock()

	if retryAfter := e.process(w.processor); retryAfter > 0 {
		w.entries.Unshift(e)

		w.mu.Lock()
		w.t.Reset(time.Duration(retryAfter) * time.Second)
		w.mu.Unlock()

		return true
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.t.Reset(w.d)

	return true
}
