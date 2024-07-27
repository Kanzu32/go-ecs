package main

import (
	"go-ecs/ecs"
	"go-ecs/test"
	"testing"
	"unsafe"
)

func BenchmarkUnsafePointersRun1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		w := test.CreateWorld(1000)
		position := w.CreateComponentPool(10)
		color := w.CreateComponentPool(10)
		velocity := w.CreateComponentPool(10)

		positionPool := w.GetPool(position)
		colorPool := w.GetPool(color)
		velocityPool := w.GetPool(velocity)

		for i := 0; i < 400; i++ {
			positionPool.AddEntity(i, unsafe.Pointer(&Position{i, i}))
		}

		for i := 250; i < 800; i++ {
			colorPool.AddEntity(i, unsafe.Pointer(&Color{"black"}))
		}

		for i := 700; i < 1000; i++ {
			velocityPool.AddEntity(i, unsafe.Pointer(&Velocity{i, i}))
		}

		arrID := make([]int, 0)
		b.StartTimer()

		for _, entityID := range w.GetEntitiesByComponent(position) {
			(*Position)(positionPool.GetComponent(entityID)).x = 100
		}

		for _, entityID := range w.GetEntitiesByComponent(velocity) {
			(*Velocity)(velocityPool.GetComponent(entityID)).y = 100
		}

		for _, entityID := range w.GetEntitiesByComponent(color) {
			(*Color)(colorPool.GetComponent(entityID)).color = "green"
		}

		for _, entityID := range w.GetEntitiesByFilter([]int{position}, []int{color}) {
			arrID = append(arrID, entityID)
		}

		for _, entityID := range w.GetEntitiesByFilter([]int{position}, []int{velocity}) {
			arrID = append(arrID, entityID)
		}
	}
}

func BenchmarkInterfacesRun1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		w := ecs.CreateWorld(1000)
		position := w.CreateComponentPool(10)
		color := w.CreateComponentPool(10)
		velocity := w.CreateComponentPool(10)

		positionPool := w.GetPool(position)
		colorPool := w.GetPool(color)
		velocityPool := w.GetPool(velocity)

		for i := 0; i < 400; i++ {
			positionPool.AddEntity(i, Position{i, i})
		}

		for i := 250; i < 800; i++ {
			colorPool.AddEntity(i, Color{"black"})
		}

		for i := 700; i < 1000; i++ {
			velocityPool.AddEntity(i, Velocity{i, i})
		}

		arrID := make([]int, 0)
		b.StartTimer()

		var pos Position
		for _, entityID := range w.GetEntitiesByComponent(position) {
			pos = positionPool.GetComponent(entityID).(Position)
			pos.x = 100
			positionPool.SetComponent(entityID, pos)
		}

		var vel Velocity
		for _, entityID := range w.GetEntitiesByComponent(velocity) {
			vel = velocityPool.GetComponent(entityID).(Velocity)
			vel.x = 100
			velocityPool.SetComponent(entityID, vel)
		}

		var col Color
		for _, entityID := range w.GetEntitiesByComponent(color) {
			col = colorPool.GetComponent(entityID).(Color)
			col.color = "green"
			colorPool.SetComponent(entityID, col)
		}

		for _, entityID := range w.GetEntitiesByFilter([]int{position}, []int{color}) {
			arrID = append(arrID, entityID)
		}

		for _, entityID := range w.GetEntitiesByFilter([]int{position}, []int{velocity}) {
			arrID = append(arrID, entityID)
		}
	}
}
