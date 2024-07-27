package ecs

import (
	"fmt"
)

// ###

type Component interface {
	// TODO INIT MAYBE???
}

// WORLD

type world struct {
	componentPools []*componentPool

	maxEntities int
}

func CreateWorld(maxEntities int) *world {
	w := world{}
	w.maxEntities = maxEntities
	w.componentPools = make([]*componentPool, 0)
	return &w
}

// func (w world) String() string {
// 	str := ""
// 	for id, v := range w.componentPools {
// 		str += fmt.Sprintf("Pool: %d\n%s\n", id, v.String())
// 	}
// 	return str
// }

func (w *world) GetPool(poolID int) *componentPool {
	return w.componentPools[poolID]
}

func (w *world) CreateComponentPool(pageSize int) int {
	pool := componentPool{}
	pool.sparseEntities = createPageArray(pageSize, w.maxEntities)
	w.componentPools = append(w.componentPools, &pool)
	return len(w.componentPools) - 1
}

func Pr(w *world) { // test
	for _, pool := range w.componentPools {
		fmt.Println(pool)
	}
}

func (w world) GetEntitiesByFilter(include []int, exclude []int) []int {
	shortestId := 0
	shortestLen := len(w.GetEntitiesByComponent(include[0]))
	res := make([]int, 0)
	for poolId := range include {
		if shortestLen > len(w.GetEntitiesByComponent(poolId)) {
			shortestLen = len(w.GetEntitiesByComponent(poolId))
			shortestId = poolId
		}
	}

EntityLoop:
	for _, entityId := range w.GetEntitiesByComponent(shortestId) {
		for _, poolId := range include {
			if !w.GetPool(poolId).HasEntity(entityId) {
				continue EntityLoop
			}
		}

		for _, poolId := range exclude {
			if w.GetPool(poolId).HasEntity(entityId) {
				continue EntityLoop
			}
		}

		res = append(res, entityId)
	}

	return res
}

func (w world) IsEntityInPool(entityId int, poolId int) bool {
	return w.GetPool(poolId).sparseEntities.get(entityId) != -1
}

func (w world) GetEntitiesByComponent(componentId int) []int {
	return w.GetPool(componentId).denseEntities
}

// PAGE ARRAY

type pageArray struct {
	data      [][]int
	pageSize  int
	arraySize int
}

func createPageArray(pageSize int, arraySize int) pageArray {
	p := pageArray{}
	p.pageSize = pageSize
	p.arraySize = arraySize
	p.data = make([][]int, arraySize/pageSize)
	return p
}

func (p pageArray) set(index int, value int) {
	pageNumber := index / p.pageSize
	pageIndex := index % p.pageSize

	if p.data[pageNumber] == nil {
		p.data[pageNumber] = make([]int, p.pageSize)

		p.data[pageNumber][0] = -1

		for j := 1; j < len(p.data[pageNumber]); j *= 2 {
			copy(p.data[pageNumber][j:], p.data[pageNumber][:j])
		}
	}

	p.data[pageNumber][pageIndex] = value
}

func (p pageArray) get(index int) int {
	pageNumber := index / p.pageSize
	pageIndex := index % p.pageSize

	if p.data[pageNumber] == nil {
		return -1
	}

	return p.data[pageNumber][pageIndex]
}

func (p pageArray) String() string {
	return fmt.Sprintf("Page size: %d\nArray size: %d\n%v", p.pageSize, p.arraySize, p.data)
}

// POOL

type componentPool struct {
	denseComponents []Component

	denseEntities []int

	sparseEntities pageArray
}

func (pool *componentPool) AddEntity(entityID int, c Component) (int, Component) {
	pool.denseEntities = append(pool.denseEntities, entityID)
	pool.denseComponents = append(pool.denseComponents, c)
	pool.sparseEntities.set(entityID, len(pool.denseEntities)-1)
	return len(pool.denseEntities) - 1, pool.denseComponents[len(pool.denseComponents)-1]
}

func (pool componentPool) HasEntity(entityID int) bool {
	return pool.sparseEntities.get(entityID) != -1
}

func (pool componentPool) GetEntities() []int {
	return pool.denseEntities
}

func (pool componentPool) GetComponent(entityID int) Component {
	return pool.denseComponents[pool.sparseEntities.get(entityID)]
}

func (pool componentPool) SetComponent(entityID int, c Component) {
	pool.denseComponents[pool.sparseEntities.get(entityID)] = c
}

func (pool componentPool) GetComponents() []Component {
	return pool.denseComponents
}

// func (pool componentPool[poolType]) GetDenseEntities() []int {
// 	return pool.denseEntities
// }

// func (pool componentPool[poolType]) GetDenseComponents() []poolType {
// 	return pool.denseComponents
// }

// func (pool componentPool[poolType]) GetComponentByEntity(entityId int) poolType {
// 	return pool.denseComponents[entityId]
// }

func (pool componentPool) EntityCount() int {
	return len(pool.denseEntities)
}

func (pool componentPool) String() string {
	return fmt.Sprintf("Components: %v\nDense ent: %v\nSparse ent:\n%v", pool.denseComponents, pool.denseEntities, pool.sparseEntities.String())
}
