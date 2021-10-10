package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

var (
	PopSize        = 500
	MutationRate   = 0.005
	ExtinctionRate = 0.01
)

type Person struct {
	DNA     []byte
	Fitness float64
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())

	target := []byte("to be or not to be")
	population := createPopulation(target)

	found := false
	generation := 0

	for !found {
		generation++
		bestPerson := getBest(population)
		fmt.Printf("\r generation:%d | %s | fitness: %2f", generation, string(bestPerson.DNA), bestPerson.Fitness)

		if bytes.Equal(bestPerson.DNA, target) {
			found = true
		} else {
			pool := createPool(population, target, bestPerson.Fitness)
			population = naturalSelection(pool, population, target)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("\nTime taken:%s\n", elapsed)
}

func createPerson(target []byte) (person Person) {
	personDNA := make([]byte, len(target))
	for i := 0; i < len(target); i++ {
		personDNA[i] = byte(rand.Intn(95) + 32)
	}
	person = Person{
		DNA:     personDNA,
		Fitness: 0,
	}
	person.calcFitness(target)
	return person
}

func (p *Person) calcFitness(target []byte) {
	score := 0
	for i := 0; i < len(target); i++ {
		if p.DNA[i] == target[i] {
			score++
		}
	}
	p.Fitness = float64(score) / float64(len(target))
}

func createPopulation(target []byte) (population []Person) {
	population = make([]Person, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = createPerson(target)
	}
	return population
}

func createPool(population []Person, target []byte, maxFitness float64) (pool []Person) {
	pool = make([]Person, 0)
	for i := 0; i < len(population); i++ {
		population[i].calcFitness(target)
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}

	if rand.Float64() < ExtinctionRate {
		i := rand.Intn(len(population))
		j := rand.Intn(len(population))
		pool = []Person{population[i], population[j]}
	}

	return pool
}

func naturalSelection(pool []Person, population []Person, target []byte) []Person {
	nextGeneration := make([]Person, len(population))

	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		child := crossover(a, b)
		child.mutate()
		child.calcFitness(target)

		nextGeneration[i] = child
	}
	return nextGeneration
}

func crossover(p1 Person, p2 Person) Person {
	child := Person{
		DNA:     make([]byte, len(p1.DNA)),
		Fitness: 0,
	}

	for i := 0; i < len(p1.DNA); i += 2 {
		child.DNA[i] = p1.DNA[i]
	}

	for j := 1; j < len(p2.DNA); j += 2 {
		child.DNA[j] = p2.DNA[j]
	}

	return child
}

func (p *Person) mutate() {
	for i := 0; i < len(p.DNA); i++ {
		if rand.Float64() < MutationRate {
			p.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}

func getBest(population []Person) Person {
	best := 0.0
	index := 0

	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}
