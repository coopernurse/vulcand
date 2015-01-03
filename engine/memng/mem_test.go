package memng

import (
	"testing"

	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/log"
	"github.com/coopernurse/vulcand/engine/test"
	"github.com/coopernurse/vulcand/plugin/registry"

	. "github.com/coopernurse/vulcand/Godeps/_workspace/src/gopkg.in/check.v1"
)

func TestMem(t *testing.T) { TestingT(t) }

type MemSuite struct {
	suite test.EngineSuite
	stopC chan bool
}

var _ = Suite(&MemSuite{})

func (s *MemSuite) SetUpSuite(c *C) {
	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})
}

func (s *MemSuite) SetUpTest(c *C) {
	engine := New(registry.GetRegistry())

	s.suite.ChangesC = make(chan interface{})
	s.stopC = make(chan bool)
	go engine.Subscribe(s.suite.ChangesC, s.stopC)
	s.suite.Engine = engine
}

func (s *MemSuite) TearDownTest(c *C) {
	close(s.stopC)
	s.suite.Engine.Close()
}

func (s *MemSuite) TestHostCRUD(c *C) {
	s.suite.HostCRUD(c)
}

func (s *MemSuite) TestHostWithKeyPair(c *C) {
	s.suite.HostWithKeyPair(c)
}

func (s *MemSuite) TestHostUpsertKeyPair(c *C) {
	s.suite.HostUpsertKeyPair(c)
}

func (s *MemSuite) TestListenerCRUD(c *C) {
	s.suite.ListenerCRUD(c)
}

func (s *MemSuite) TestBackendCRUD(c *C) {
	s.suite.BackendCRUD(c)
}

func (s *MemSuite) TestBackendDeleteUsed(c *C) {
	s.suite.BackendDeleteUsed(c)
}

func (s *MemSuite) TestServerCRUD(c *C) {
	s.suite.ServerCRUD(c)
}

func (s *MemSuite) TestFrontendCRUD(c *C) {
	s.suite.FrontendCRUD(c)
}

func (s *MemSuite) TestFrontendBadBackend(c *C) {
	s.suite.FrontendBadBackend(c)
}

func (s *MemSuite) TestMiddlewareCRUD(c *C) {
	s.suite.MiddlewareCRUD(c)
}

func (s *MemSuite) TestMiddlewareBadFrontend(c *C) {
	s.suite.MiddlewareBadFrontend(c)
}

func (s *MemSuite) TestMiddlewareBadType(c *C) {
	s.suite.MiddlewareBadType(c)
}
