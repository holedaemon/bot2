package web

import (
	"sync"

	"github.com/holedaemon/bot2/internal/web/templates"
)

type cachedGuilds struct {
	mu     sync.Mutex
	guilds map[string]*templates.Guild
}

func newGuildCache() *cachedGuilds {
	return &cachedGuilds{
		guilds: make(map[string]*templates.Guild),
	}
}

func (gc *cachedGuilds) Add(guild *templates.Guild) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	if gc.guilds == nil {
		gc.guilds = make(map[string]*templates.Guild)
	}

	gc.guilds[guild.ID] = guild
}

func (gc *cachedGuilds) Get(id string) (*templates.Guild, bool) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	if gc.guilds == nil {
		gc.guilds = make(map[string]*templates.Guild)
	}

	g, ok := gc.guilds[id]
	return g, ok
}

func (gc *cachedGuilds) ToSlice() []*templates.Guild {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	s := make([]*templates.Guild, 0, len(gc.guilds))
	for _, g := range gc.guilds {
		s = append(s, g)
	}

	return s
}
