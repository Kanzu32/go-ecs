package main

import (
	"go-ecs/ecs"
	"testing"
)

func BenchmarkUnsafePointersNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		w := ecs.CreateWorld(10000)
		positionPool, _ := ecs.CreateComponentPool[Position](w, 100)
		colorPool, _ := ecs.CreateComponentPool[Color](w, 100)
		velocityPool, _ := ecs.CreateComponentPool[Velocity](w, 100)

		// positionPool := ecs.Pool[Position](w, positionID)
		// colorPool := ecs.Pool[Color](w, colorID)
		// velocityPool := ecs.Pool[Velocity](w, velocityID)

		// fmt.Print("HERE")
		for i := 0; i < 400; i++ {
			positionPool.AddEntity(i, Position{i, i})
		}

		for i := 250; i < 800; i++ {
			colorPool.AddEntity(i, Color{"red"})
		}

		for i := 700; i < 1000; i++ {
			velocityPool.AddEntity(i, Velocity{i, i})
		}

		b.StartTimer()

		for _, entID := range velocityPool.Entities() {
			velocityPool.Component(entID).y = 100
		}

		for _, entID := range positionPool.Entities() {
			positionPool.Component(entID).x = 100
		}

		for _, entID := range colorPool.Entities() {
			colorPool.Component(entID).color = "green"
		}
	}
}

func BenchmarkUnsafePointersNewFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		w := ecs.CreateWorld(10000)
		positionPool, _ := ecs.CreateComponentPool[Position](w, 100)
		colorPool, _ := ecs.CreateComponentPool[Color](w, 100)
		velocityPool, _ := ecs.CreateComponentPool[Velocity](w, 100)

		// positionPool := ecs.Pool[Position](w, positionID)
		// colorPool := ecs.Pool[Color](w, colorID)
		// velocityPool := ecs.Pool[Velocity](w, velocityID)

		// fmt.Print("HERE")
		for i := 0; i < 400; i++ {
			positionPool.AddEntity(i, Position{i, i})
		}

		for i := 250; i < 800; i++ {
			colorPool.AddEntity(i, Color{"red"})
		}

		for i := 700; i < 1000; i++ {
			velocityPool.AddEntity(i, Velocity{i, i})
		}

		arrID := make([]int, 0)
		b.StartTimer()

		for _, entID := range ecs.PoolFilter([]ecs.P{positionPool, colorPool}, []ecs.P{}) {
			arrID = append(arrID, entID)
		}

		for _, entID := range ecs.PoolFilter([]ecs.P{velocityPool, colorPool}, []ecs.P{}) {
			arrID = append(arrID, entID)
		}
	}
}
