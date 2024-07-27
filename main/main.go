package main

import (
	"fmt"
	"go-ecs/ecs"
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

type Struct struct {
	arr []ecs.Component
}

func main() {
	w := ecs.CreateWorld(1000)
	position := w.CreateComponentPool(10)

	positionPool := w.GetPool(position)
	positionPool.AddEntity(0, Position{2, 3})

	pos := positionPool.GetComponent(0).(Position)
	pos.x = 100
	positionPool.SetComponent(0, pos)
	fmt.Println(positionPool.GetComponent(0).(Position).x)
}
