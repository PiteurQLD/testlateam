package scrapper

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type User struct {
	UID           string
	Name          string
	mtx           sync.RWMutex
	apiBase       string
	DelayInterval time.Duration
	headers       map[string]string
	positions     map[string]Position
	client        *http.Client
	log           *log.Logger
	firstFetch    bool
}
type UserOption func(*User)

func NewUser(UID string, Name string, DelayInterval time.Duration, opts ...UserOption) *User {
	u := User{
		UID:           UID,
		Name:          Name,
		log:           logger,
		positions:     make(map[string]Position),
		DelayInterval: DelayInterval,
		client:        http.DefaultClient,
		firstFetch:    true,
		headers:       defaultHeaders,
		apiBase:       defaultApiBase,
	}
	u.log.SetOutput(io.Discard)
	for _, opt := range opts {
		opt(&u)
	}
	return &u
}
func (u *User) SetAPIBase(s string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.apiBase = s
}
func (u *User) APIBase() string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.apiBase
}
func (u *User) SetDelay(d time.Duration) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.DelayInterval = d
}
func (u *User) Delay() time.Duration {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.DelayInterval
}
func (u *User) SetHeaders(h map[string]string) {
	headers := make(map[string]string, len(h))
	for k, v := range h {
		headers[k] = v
	}
	u.mtx.Lock()
	defer u.mtx.Unlock()
	u.headers = h
}
func (u *User) Headers() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	headers := make(map[string]string, len(u.headers))
	for k, v := range u.headers {
		headers[k] = v
	}
	return headers
}
func WithCustomLogger(l *log.Logger) UserOption {
	return func(u *User) {
		u.log = l
	}
}
func WithLogging() UserOption {
	return func(u *User) {
		u.log.SetOutput(os.Stdout)
	}
}
func WithCustomRefresh(d time.Duration) UserOption {
	return func(u *User) {
		u.DelayInterval = d
	}
}
func WithHTTPClient(c *http.Client) UserOption {
	return func(u *User) {
		u.client = c
	}
}
func WithHeaders(h map[string]string) UserOption {
	return func(u *User) {
		u.headers = h
	}
}
func WithTestnet() UserOption {
	return func(u *User) {
		u.SetAPIBase("https://testnet.binancefuture.com/bapi/future")
	}
}
