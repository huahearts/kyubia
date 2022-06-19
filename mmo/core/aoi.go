package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

type AOIManager struct {
	MinX  int
	MaxX  int
	CntsX int
	MinY  int
	MaxY  int
	CntsY int
	Grids map[int]*Grid
}

func NewAOIMgr(minx, maxx, cntsX, miny, maxy, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minx,
		MaxX:  maxx,
		MinY:  miny,
		MaxY:  maxy,
		CntsX: cntsX,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}
	for y := 0; y < aoiMgr.CntsY; y++ {
		for x := 0; x < aoiMgr.CntsX; x++ {
			gid := y*aoiMgr.CntsX + x
			aoiMgr.Grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}
	return aoiMgr
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

func (m *AOIManager) GetSurroundGridByGid(gid int) (grids []*Grid) {
	if _, ok := m.Grids[gid]; !ok {
		return
	}

	grids = append(grids, m.Grids[gid])

	idx := gid % m.CntsX

	if idx > 0 {
		grids = append(grids, m.Grids[gid-1])
	}

	if idx < m.CntsX-1 {
		grids = append(grids, m.Grids[gid+1])
	}

	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几列
		idy := v / m.CntsX

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CntsX])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.Grids[v+m.CntsX])
		}
	}

	return
}

//通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(x) - m.MinY) / m.gridLength()

	return gy*m.CntsX + gx
}

//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundGridByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		fmt.Printf("===> grid ID : %d, pids : %v  ====", v.GID, v.GetPlayerIDs())
	}

	return
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetPlayerIDs()
	return
}

//移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].Add(pID)
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Remove(pID)
}
