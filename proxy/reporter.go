package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/oxy/memmetrics"
	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/oxy/utils"
	"github.com/coopernurse/vulcand/engine"

	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/metrics"
	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/timetools"
)

// Reporter reports real time metrics to the Statsd client
type Reporter struct {
	next  http.Handler
	c     metrics.Client
	fk    engine.FrontendKey
	clock timetools.TimeProvider
}

func NewReporter(next http.Handler, c metrics.Client, fk engine.FrontendKey) *Reporter {
	return &Reporter{
		clock: &timetools.RealTime{},
		next:  next,
		c:     c,
		fk:    fk,
	}
}

func (rp *Reporter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := rp.clock.UtcNow()
	pw := &utils.ProxyWriter{W: w}
	rp.next.ServeHTTP(pw, req)
	diff := rp.clock.UtcNow().Sub(start)
	rp.emitMetrics(req, pw.Code, diff, "frontend", rp.fk.Id)
}

func (rp *Reporter) emitMetrics(r *http.Request, code int, latency time.Duration, p ...string) {
	m := rp.c.Metric(p...)
	rp.c.TimingMs(m.Metric("rtt"), latency, 1)
	rp.c.Inc(m.Metric("code", fmt.Sprintf("%v", code)), 1, 1)
	rp.c.Inc(m.Metric("request"), 1, 1)
}

// RTWatcher watches runtime metrics
type RTWatcher struct {
	mtx   *sync.Mutex
	m     *memmetrics.RTMetrics
	srvs  map[surl]*memmetrics.RTMetrics
	clock timetools.TimeProvider
	next  http.Handler
}

func NewWatcher(next http.Handler) (*RTWatcher, error) {
	m, err := memmetrics.NewRTMetrics()
	if err != nil {
		return nil, err
	}

	return &RTWatcher{
		mtx:   &sync.Mutex{},
		m:     m,
		clock: &timetools.RealTime{},
		next:  next,
		srvs:  make(map[surl]*memmetrics.RTMetrics),
	}, nil
}

func (rt *RTWatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := rt.clock.UtcNow()
	pw := &utils.ProxyWriter{W: w}
	rt.next.ServeHTTP(pw, req)
	diff := rt.clock.UtcNow().Sub(start)

	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	rt.m.Record(pw.Code, diff)

	sm, ok := rt.srvs[surl{scheme: req.URL.Scheme, host: req.URL.Host}]
	if ok {
		sm.Record(pw.Code, diff)
	}
}

func (rt *RTWatcher) rtStats() (*engine.RoundTripStats, error) {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	return engine.NewRoundTripStats(rt.m)
}

func (rt *RTWatcher) rtServerStats(u *url.URL) (*engine.RoundTripStats, error) {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	sm, ok := rt.srvs[surl{scheme: u.Scheme, host: u.Host}]
	if ok {
		return engine.NewRoundTripStats(sm)
	}
	return nil, fmt.Errorf("watcher: %v not found", u)
}

func (rt *RTWatcher) upsertServer(u *url.URL) error {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	m, err := memmetrics.NewRTMetrics()
	if err != nil {
		return err
	}
	rt.srvs[surl{scheme: u.Scheme, host: u.Host}] = m
	return nil
}

func (rt *RTWatcher) hasServer(u *url.URL) bool {
	_, ok := rt.srvs[surl{scheme: u.Scheme, host: u.Host}]
	return ok
}

func (rt *RTWatcher) removeServer(u *url.URL) {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	delete(rt.srvs, surl{scheme: u.Scheme, host: u.Host})
}

func (rt *RTWatcher) collectMetrics(m *memmetrics.RTMetrics) error {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	return m.Append(rt.m)
}

func (rt *RTWatcher) collectServerMetrics(m *memmetrics.RTMetrics, u *url.URL) error {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	sm, ok := rt.srvs[surl{scheme: u.Scheme, host: u.Host}]
	if ok {
		m.Append(sm)
	}
	return nil
}

type surl struct {
	scheme string
	host   string
}
