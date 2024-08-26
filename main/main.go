package main

import (
	"fmt"
	"go-ecs/ecs"
	"go-ecs/ecs/psize"
)

type Position struct {
	x int
	y int
}

type Color struct {
	color string
}

type Velocity struct {
	x int
	y int
}

func main() {
	w := ecs.CreateWorld()
	positionPool := ecs.CreateComponentPool[Position](w, psize.Page128)
	// colorPool := ecs.CreateComponentPool[Color](w, 100)
	// moveFlag := ecs.CreateFlagPool(w, 10)

	// positionPool.AddNewEntity(Position{0, 0})
	// positionPool.AddNewEntity(Position{1, 1})
	// ent, _ := positionPool.AddNewEntity(Position{2, 2})
	// positionPool.AddNewEntity(Position{3, 3})
	// positionPool.AddNewEntity(Position{4, 4})
	// fmt.Println(positionPool)
	// ecs.RemoveEntityFromWorld(w, ent)
	// fmt.Println(positionPool)

	for i := 0; i < 6; i++ {
		positionPool.AddNewEntity(Position{i, i})
	}
	// for _, entity := range positionPool.Entities() {
	// 	if entity.ID() >= 300 && entity.ID() < 700 {
	// 		colorPool.AddExistingEntity(entity, Color{"black"})
	// 	}
	// 	if entity.ID() >= 500 && entity.ID() < 1000 {
	// 		moveFlag.AddExistingEntity(entity)
	// 	}
	// 	pos, _ := positionPool.Component(entity)
	// 	pos.y++
	// }
	var ent ecs.Entity

	for _, entity := range positionPool.Entities() {
		if entity.ID() == 2 {
			ent = entity
		}
	}

	for _, entity := range positionPool.Entities() {
		if entity.ID() == 2 {
			ecs.RemoveEntityFromWorld(w, entity)
		}
	}

	newent, _ := positionPool.AddNewEntity(Position{-1, -1})

	for _, entity := range ecs.PoolFilter([]ecs.AnyPool{positionPool}, []ecs.AnyPool{}) {
		comp, _ := positionPool.Component(entity)
		fmt.Println(entity, comp.x, comp.y)
	}

	fmt.Println(ent)
	fmt.Println(newent)
	fmt.Println(newent == ent)

	// for i := 100; i < 300; i++ {
	// 	ecs.RemoveEntityFromWorld(w, i)
	// }

	// for _, entityID := range ecs.PoolFilter([]ecs.AnyPool{moveFlag}, []ecs.AnyPool{}) {
	// 	fmt.Println(entityID, "flag")
	// }

}
