package aoi

import "errors"

type grid struct {
	r, c     int32                       // row/col coordinate.
	entities map[interface{}]*GridEntity // AOI entities inside.
}

// Grid AOI manager.
type GridManager struct {
	lb, rt Position // left-bottom, right-top.
	gWidth int32    // grid width.
	grids  [][]grid // all grids.
}

func NewGrid(lb, rt Position, gWidth int32) *GridManager {
	if gWidth <= 0 {
		panic(errors.New("invalid width"))
	}

	width, height := rt.X-lb.X, rt.Y-lb.Y
	r, c := width/gWidth, height/gWidth
	if width%gWidth > 0 {
		r += 1
	}
	if height%gWidth > 0 {
		c += 1
	}

	g := &GridManager{
		lb:     lb,
		rt:     rt,
		gWidth: gWidth,
	}
	g.grids = make([][]grid, r)
	for i := int32(0); i < r; i++ {
		g.grids[i] = make([]grid, c)
		for j := int32(0); j < c; j++ {
			g.grids[i][j] = grid{
				r:        i,
				c:        j,
				entities: map[interface{}]*GridEntity{},
			}
		}
	}

	return g
}

func (m *GridManager) Add(key interface{}, pos Position, user User) (Entity, error) {
	grid := m.getGridByPos(pos)
	if grid == nil {
		return nil, errors.New("invalid position")
	}

	if user == nil {
		return nil, errors.New("nil user")
	}

	if _, ok := grid.entities[key]; ok {
		return nil, errors.New("entity already exist")
	}

	entity := &GridEntity{
		entity:  newEntity(key, pos, user),
		manager: m,
	}

	var entities []User
	enter := []User{user}
	grids := m.getNearbyGrids(grid)
	for _, g := range grids {
		for _, e := range g.entities {
			entities = append(entities, e.user)
			e.user.OnAOIUpdate(enter, nil)
		}
	}

	grid.entities[key] = entity
	entity.grid = grid
	entity.user.OnAOIUpdate(entities, nil)

	return entity, nil
}

func (m *GridManager) Rem(e Entity) error {
	ee, ok := e.(*GridEntity)
	if !ok {
		return errors.New("invalid entity")
	}

	if ee.manager != m {
		return errors.New("object not belongs to manager")
	}

	grid := ee.grid
	delete(ee.grid.entities, ee.key)

	leave := []User{e.User()}
	grids := m.getNearbyGrids(grid)
	for _, g := range grids {
		for _, v := range g.entities {
			v.user.OnAOIUpdate(nil, leave)
		}
	}

	ee.manager = nil
	ee.grid = nil
	return nil
}

func (m *GridManager) PosNearAOI(pos Position, distance int32) []User {

	lbp := m.fixPos(Position{
		X: pos.X - distance,
		Y: pos.Y - distance,
	})
	rtp := m.fixPos(Position{
		X: pos.X + distance,
		Y: pos.Y + distance,
	})

	lbg := m.getGridByPos(lbp)
	rtg := m.getGridByPos(rtp)

	grids := []*grid{}
	for r := lbg.r; r <= rtg.r; r++ {
		for c := lbg.c; c <= rtg.c; c++ {
			grids = append(grids, &m.grids[r][c])
		}
	}

	ret := make([]User, 0)
	for _, g := range grids {
		for _, e := range g.entities {
			ret = append(ret, e.user)
		}
	}
	return ret
}

func (m *GridManager) fixPos(pos Position) Position {
	if pos.X < m.lb.X {
		pos.X = m.lb.X
	}
	if pos.X > m.rt.X {
		pos.X = m.rt.X
	}
	if pos.Y < m.lb.X {
		pos.Y = m.lb.X
	}
	if pos.Y > m.rt.Y {
		pos.Y = m.rt.Y
	}
	return pos
}

func (m *GridManager) getGridByPos(pos Position) *grid {
	if pos.X < m.lb.X || pos.X > m.rt.X || pos.Y < m.lb.Y || pos.Y > m.rt.Y {
		return nil
	}

	r := (pos.X - m.lb.X) / m.gWidth
	c := (pos.Y - m.lb.Y) / m.gWidth
	if r == int32(len(m.grids)) {
		r -= 1
	}
	if c == int32(len(m.grids[0])) {
		c -= 1
	}
	return &m.grids[r][c]
}

func (m *GridManager) getNearbyGrids(g *grid) []*grid {
	grids := make([]*grid, 0, 9)
	rc := [9]struct{ r, c int32 }{
		{g.r - 1, g.c - 1}, {g.r - 1, g.c}, {g.r - 1, g.c + 1},
		{g.r, g.c - 1}, {g.r, g.c}, {g.r, g.c + 1},
		{g.r + 1, g.c - 1}, {g.r + 1, g.c}, {g.r + 1, g.c + 1},
	}

	maxR, maxC := int32(len(m.grids))-1, int32(len(m.grids[0]))-1

	for i := 0; i < 9; i++ {
		if rc[i].r < 0 || rc[i].r > maxR || rc[i].c < 0 || rc[i].c > maxC {
			continue
		}
		grids = append(grids, &m.grids[rc[i].r][rc[i].c])
	}

	return grids
}

type GridEntity struct {
	entity
	manager *GridManager
	grid    *grid
}

func (e *GridEntity) Move(pos Position) error {
	if pos == e.pos {
		return nil
	}

	newGrid := e.manager.getGridByPos(pos)
	if newGrid == nil {
		return errors.New("invalid position")
	}

	oldGrid := e.grid
	if oldGrid.r == newGrid.r && oldGrid.c == newGrid.c {
		return nil
	}

	oldGrids := e.manager.getNearbyGrids(oldGrid)
	newGrids := e.manager.getNearbyGrids(newGrid)
	oldLast := len(oldGrids) - 1
	newLast := len(newGrids) - 1
	delete(oldGrid.entities, e.key)

	selfLeave := []User{e.user}
	othLeave := make([]User, 0)
	for _, g := range oldGrids {
		if !isGridInsideGrids(newGrids[0], newGrids[newLast], g) {
			for _, v := range g.entities {
				v.user.OnAOIUpdate(nil, selfLeave)
				othLeave = append(othLeave, v.user)
			}
		}
	}

	selfEnter := []User{e.user}
	othEnter := make([]User, 0)
	for _, g := range newGrids {
		if !isGridInsideGrids(oldGrids[0], oldGrids[oldLast], g) {
			for _, v := range g.entities {
				v.user.OnAOIUpdate(selfEnter, nil)
				othEnter = append(othEnter, v.user)
			}
		}
	}

	newGrid.entities[e.key] = e
	e.grid = newGrid
	e.pos = pos
	e.user.OnAOIUpdate(othEnter, othLeave)

	return nil
}

func (e *GridEntity) TraverseAOI(fn func(u User) error) error {
	if fn == nil {
		return errors.New("nil fn")
	}

	grids := e.manager.getNearbyGrids(e.grid)
	for _, g := range grids {
		if e.grid != g {
			for _, v := range g.entities {
				if err := fn(v.user); err != nil {
					return err
				}
			}
		} else {
			for _, v := range g.entities {
				if e != v {
					if err := fn(v.user); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func isGridInsideGrids(leftTop, rightBottom, g *grid) bool {
	return g.r >= leftTop.r && g.r <= rightBottom.r && g.c >= leftTop.c && g.c <= rightBottom.c
}
