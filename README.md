## [EN] Description
Implementation of the Enity-Component-System architectural pattern in Go.

The ECS pattern is widely used in game development, high-performance simulations, and other systems where performance and flexibility are important.

This library uses the sparse-set approach, which allows to perform operations on entities in O(1) time. Cache efficiency is also improved by storing data of the same type close to each other in memory.

The game [Troublemakers](https://github.com/Kanzu32/strategy-game) was created based on this library.

## Features
* Versions of objects for efficient disposal;
* Sparse-set approach.

## Technologies
* Golang;
* Data-oriented;
* Entity-Component-System.

## [RU] Описание
Реализация архитектурного паттерна Enity-Component-System на языке Go. 

Паттерн ECS широко применяется при разработке игр, высокопроизводительных симуляциях и других cистемах, где важна производительность и гибкость.

Данная библиотека использует sparse-set подход, позволяющий выполнять операции над сущностями за время O(1). Также повышается кэш-эффективность из-за хранения данных одного типа вплотную друг к другу в памяти.

На основе данной библиотеки создана игра [Смутьяны](https://github.com/Kanzu32/strategy-game).

## Особенности
* Версии объектов для эффективной утилизации;
* Sparse-set подход.

## Модель данных пула объектов
![](https://github.com/Kanzu32/go-ecs/blob/main/readme/go-ecs-1.png)
