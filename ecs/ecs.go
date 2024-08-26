package ecs

import (
	"errors"
	"fmt"
	"go-ecs/ecs/parray"
	"go-ecs/ecs/psize"
)

const maxEntities = 65536

type P = AnyPool
type AnyPool interface {
	HasEntity(entity Entity) bool
	Entities() []Entity
	EntityCount() int
	RemoveEntity(entity Entity) error
}

type Entity struct {
	state   uint8
	id      uint16
	version uint8
}

func (e *Entity) ID() uint16         { return e.id }
func (e *Entity) isNil() bool        { return e.state&1 != 0 }
func (e *Entity) isRegistered() bool { return e.state&2 != 0 }
func (e *Entity) setNil() {
	e.state = e.state | 1
	e.state = e.state &^ 2
}
func (e *Entity) setRegistered() {
	e.state = e.state | 2
	e.state = e.state &^ 1
}

func (e *Entity) clear() { e.state = 0 }
func (e Entity) String() string {
	if e.isNil() {
		return fmt.Sprintf("E #%d v%d NIL", e.id, e.version)
	}
	if e.isRegistered() {
		return fmt.Sprintf("E #%d v%d REG", e.id, e.version)
	}
	return fmt.Sprintf("Ent #%d v%d ", e.id, e.version)
}

// func (f *flags) Set(b)      { return b | flag }
// func (f *flags) Clear(b)    { return b &^ flag }
// func (f *flags) Toggle(b)   { return b ^ flag }
// func (f *flags) Has(b) bool { return b&flag != 0 }

// WORLD

type world struct {
	pools []AnyPool

	next      uint32   // next available entity ID
	entities  []Entity // array to mark registred and destroyed entities
	destroyed Entity   // last entity removed from forld
}

func CreateWorld() *world {
	w := world{}
	w.pools = make([]AnyPool, 0)
	w.entities = make([]Entity, maxEntities)
	e := Entity{}
	e.setNil()
	w.destroyed = e
	return &w
}

func CreateComponentPool[componentType any](w *world, pageSize psize.PageSizes) *componentPool[componentType] {
	pool := componentPool[componentType]{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func CreateFlagPool(w *world, pageSize psize.PageSizes) *flagPool {
	pool := flagPool{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func (w *world) registerNewEntity() (Entity, error) {
	if w.destroyed.isNil() {
		if w.next == maxEntities {
			e := Entity{}
			e.setNil()
			return e, errors.New("too many entities")
		}
		w.entities[w.next].setRegistered()
		e := w.entities[w.next]
		e.id = uint16(w.next)
		w.next += 1
		return e, nil
	}

	ret := w.destroyed
	w.destroyed = w.entities[ret.id]
	w.entities[ret.id].setRegistered()
	ret.setRegistered()
	ret.version = w.entities[ret.id].version
	return ret, nil
}

func (w *world) isRegisteredEntity(entity Entity) bool {
	return w.entities[entity.id].isRegistered()
}

func RemoveEntityFromWorld(w *world, entity Entity) {
	for _, pool := range w.pools {
		if pool.HasEntity(entity) {
			pool.RemoveEntity(entity)
		}
	}
	w.entities[entity.id].id = w.destroyed.id
	w.entities[entity.id].state = w.destroyed.state
	w.entities[entity.id].version++
	w.destroyed.id = entity.id
	w.destroyed.clear()
}

// COMPONENT POOL

type componentPool[componentType any] struct {
	denseComponents []componentType

	denseEntities []Entity

	sparseEntities parray.PageArray

	world *world
}

func (pool *componentPool[componentType]) AddNewEntity(comp componentType) (Entity, error) {
	entity, err := pool.world.registerNewEntity()
	if err != nil {
		return entity, err
	}
	pool.denseComponents = append(pool.denseComponents, comp)
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *componentPool[componentType]) AddExistingEntity(entity Entity, comp componentType) error {
	if !pool.world.isRegisteredEntity(entity) {
		return errors.New("entityID is not registered")
	}
	pool.denseComponents = append(pool.denseComponents, comp)
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return nil
}

func (pool *componentPool[componentType]) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.id)                                     // индекс для удаления (замены) элемента в dense массивах
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].id                        // индекс элемента в sparse массиве для последнего dense элемента
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)                                 // установка нового указателя на dence массив sparce массиве
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]       // перемещение последнего элемента dense массива
	pool.denseComponents[denseRemoveIndex] = pool.denseComponents[len(pool.denseComponents)-1] // на позицию удаления для двух массивов

	pool.sparseEntities.Set(entity.id, -1) // установка sparse эдемента для удаления в -1

	// Уменьшение len без удаления последнего элемента.
	// При необходимости его можно восстановить увиличив len. Append перезапишет скрытый элемент
	pool.denseComponents = pool.denseComponents[:len(pool.denseComponents)-1]
	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *componentPool[componentType]) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.id) != -1
}

func (pool *componentPool[componentType]) Entities() []Entity {
	return pool.denseEntities
}

func (pool *componentPool[componentType]) Component(entity Entity) (*componentType, error) {
	return &pool.denseComponents[pool.sparseEntities.Get(entity.id)], nil
}

func (pool *componentPool[poolType]) EntityCount() int {
	return len(pool.denseEntities)
}

func (pool componentPool[componentType]) String() string {
	return fmt.Sprintf("Components: %v\nDense ent: %v\nSparse ent:\n%v", pool.denseComponents, pool.denseEntities, pool.sparseEntities.String())
}

// ENTITY FILTER

func PoolFilter(include []AnyPool, exclude []AnyPool) []Entity {
	if len(include) == 0 {
		panic("include can't be empty")
	}
	shortestIndex := 0
	shortestLen := include[0].EntityCount()
	res := make([]Entity, 0)
	for i, pool := range include {
		if shortestLen > pool.EntityCount() {
			shortestLen = pool.EntityCount()
			shortestIndex = i
		}
	}
EntityLoop:
	for _, entity := range include[shortestIndex].Entities() {
		for _, pool := range include {
			if !pool.HasEntity(entity) {
				continue EntityLoop
			}
		}

		for _, pool := range exclude {
			if pool.HasEntity(entity) {
				continue EntityLoop
			}
		}

		res = append(res, entity)
	}
	return res
}

// FLAG POOL

type flagPool struct {
	denseEntities []Entity

	sparseEntities parray.PageArray

	world *world
}

func (pool *flagPool) AddNewEntity() (Entity, error) {
	entity, err := pool.world.registerNewEntity()
	if err != nil {
		return entity, err
	}
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *flagPool) AddExistingEntity(entity Entity) (Entity, error) {
	if !pool.world.isRegisteredEntity(entity) {
		return entity, errors.New("entityID is not registered")
	}
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *flagPool) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.id)
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].id
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]

	pool.sparseEntities.Set(entity.id, -1)

	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *flagPool) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.id) != -1
}

func (pool *flagPool) Entities() []Entity {
	return pool.denseEntities
}

func (pool *flagPool) EntityCount() int {
	return len(pool.denseEntities)
}
