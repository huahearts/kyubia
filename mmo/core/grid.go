package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int
	MinX      int
	MaxX      int // 是否有必要存在
	MinY      int
	MaxY      int
	playerIDs map[int]bool
	pIDLok    sync.RWMutex
}

func NewGrid(gid, minx, maxx, miny, maxy int) *Grid {
	return &Grid{
		GID:       gid,
		MinX:      minx,
		MinY:      miny,
		MaxX:      maxx,
		MaxY:      maxy,
		playerIDs: make(map[int]bool),
	}
}

func (g *Grid) Add(playerId int) {
	g.pIDLok.Lock()
	defer g.pIDLok.Unlock()

	g.playerIDs[playerId] = true
}

func (g *Grid) Remove(playerID int) {
	g.pIDLok.Lock()
	defer g.pIDLok.Unlock()
	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIds []int) {
	g.pIDLok.Lock()
	defer g.pIDLok.Unlock()
	for k, _ := range g.playerIDs {
		playerIds = append(playerIds, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
