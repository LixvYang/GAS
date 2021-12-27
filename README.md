# GAS# 遗传算法

## 简介

遗传算法(Genetic Algorithm)是计算数学中用于解决最优化的搜索算法，基于自然选择过程。

自然选择是进化的一个关键机制。在一个自然过程中，种群会随着时间的推移而**不得不**适应环境，这个种群在性状上有差异。最终具有更合适性状的个体生物在环境中生存的机会更高。从而这些幸存下来的生物繁殖下一代来继承他们的性状，并最终产生出更合适性状的物种。但是，如果整个种群具有相同的性状，随着时间的推移，环境发生突变的时候，该物种就会灭绝。

幸运的是，遗传算法中最重要的是突变，**突变**导致的性状发生变化从而使得更适应环境变化的生物能够生存下来，这使得更适合环境的性状得以保留。

高中的时候，生物课本上有这么一张图：桦尺蠖

 [桦尺蠖.jfif](桦尺蠖.jfif) 

十九世纪之前，桦尺蠖主要是白色品种，白色可以帮助桦尺蠖免于天敌的威胁，因为它的白色与树木很好地融合在一起。然而英国工业革命期间，许多树木被熏黑，这使得深色的桦尺蠖在躲避捕食者的过程中具有优势。如今，随着保护环境活动的进行，这种平衡再一次被打破。

这就是自然选择——适者生存。遗传算法有很多与自然选择类似的机制——DNA，种群，变异，适应度，选择，繁殖，遗传和突变。

- DNA——定义生物体的遗传物质
- 种群——从最初的生物体种群开始，它们的DNA有不同的值
- 适应度——确定每个生物体对环境是适应程度
- 选择——选择适应度最好的生物，给它们更高的繁殖机会
- 繁殖——从自然选择中选定适应度最好的生物以繁殖下一代
- 继承——下一代必须继承DNA的值
- 突变——每一代中，基因有很小概率会发生突变

## 无限猴子定理

![](https://natureofcode.com/book/imgs/chapter09/ch09_01.png)

**无限猴子定理**的表述如下：让一只猴子在打字机上随机地按键，当按键时间达到无穷时，几乎必然能够打出任何给定的文字，比如莎士比亚的全套著作。这个理论的问题在于，这只猴子实际上打字莎士比亚的可能性非常低，即使那只猴子从宇宙大爆炸开始，算到现在，也不太可能拥有莎士比亚全套著作。

换个角度，让我们的猴子打印这一句话“to be or not to be that is the question”(生存还是死亡，这是一个问题)。这句话的字符为39个，如果放在猴子面前的键盘有27个键(26个英文字母和一个空格)，那么他打对第一个字母的概率为1/27，类推，他有(1/27)^39倍的概率能够一次性输入正确。

显然使用遍历的方法不会让我们有所进展，但让我们尝试遗传算法怎么样？

> 如果你已经读到了这里，请坚持读下去，因为那些阅读能力不够五行的人已经被**淘汰**了——**自然选择**

## 构思

事实上知道了以上的知识并不足以让我们完整地写出整个程序，或许在写代码之前进行完整的构思，在接下来的过程中完成细节的补充很重要。

先来想一想我们需要什么:

- 目标，即最佳的遗传物质DNA
- 个体——种群的最小组成单位
- 种群——个体的组成集合
- 适应度——淘汰种群内部不适应的个体
- 选择机会——选择适应度最好的生物以繁殖下一代
- 继承——下一代继承上一代的DNA
- 突变——为形成最好的DNA做准备
- 灭绝——小概率下，物种有灭绝风险
- ......

## 代码

#### 对象

我们定义一个个体，此对象有DNA、和适应度两种属性，分别为字节数组和代表适应度的数字

```go
type Person struct {
	DNA     []byte
	Fitness float64
}
```

没什么特别的

#### 定义种群

我们需要为种群创造生物体，代码放在这里

```go
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
```

`target`就是我们要实现的最佳DNA，在这个函数种，我们创建一个与`target`长度相同的字节数组，并为其设置新创建的个体的基因值。

计算基因值适应度的函数在这儿：

```go
func (p *Person) calcFitness(target []byte) {
	score := 0
	for i := 0; i < len(target); i++ {
		if p.DNA[i] == target[i] {
			score++
		}
	}
	p.Fitness = float64(score) / float64(len(target))
}
```

这个函数比较个体`DNA`与`target`字节匹配的次数，score除以目标的总字节数使适应度为百分比，得到`Fitness`，这是介于0.0到1.0的数字。

如果适应度为1.0，则我们个体的DNA为最佳DNA。

既然有了完整的个体，那我们就可以创建种群。

```go
func createPopulation(target []byte) (population []Person) {
	population = make([]Person, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = createPerson(target)
	}
	return population
}
```

这里的种群`population`是一个`Person`数组，`PopSize`是赋予种群大小的全局变量。

#### 选择适应度最高的生物，给予繁殖机会

现在我们已经有一个种群，找出其中适应度最好的个体以繁殖下一代，我们这里选择“[繁殖池](https://en.wikipedia.org/wiki/Mating_pool)”机制。

```go
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
```

我们所做的就是创造一个繁殖池，里面有同一个生物的多个副本放入其中。生物适应度越高，最终放入池中的生物体副本越多。

我在这里放了一个ExtinctionRate机制，就是灭绝率，如果随机数小于`ExtinctionRate`我们的繁殖池就在种群中随机选择两个个体(在灭绝的危机下，适应度就不适用了，随机性才适用)

#### 选定最适合的生物体创建下一代种群

我们已经有了一个繁殖池，并且在没有灭绝的情况下，在繁殖池中随机选择两个生物，并将它们作为亲本，为种群创造下一代。

```go
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
```

下一代必须继承上一代的DNA

```go
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
```

我们的继承方式选择交叉进行，即每个DNA的字节从两个亲本之间交替复制。

#### 随机变异

在从两个亲本个体中繁殖出新的子生物体后，我们观察生物体是否会突变。

```go
func (p *Person) mutate() {
	for i := 0; i < len(p.DNA); i++ {
		if rand.Float64() < MutationRate {
			p.DNA[i] = byte(rand.Intn(95) + 32)
		}
	}
}
```

这里的突变就是随机数小于突变率`MutationRate`，突变很重要，如果没有突变发生，那么种群中的DNA永远会与原始种群相同——永远没有最佳结果。如果没有突变，自然选择根本就不起效果。

突变可以使我们摆脱局部最大值找到全局最大值。

一旦我们发现了突变，我们就计算子个体的适应度并插入到下一代种群中。

这就是遗传算法的全部内容。

我们定义一下全局变量

```go
var (
	PopSize        = 500
	MutationRate   = 0.005
	ExtinctionRate = 0.01
)
```

main函数：

```go
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
```

我把目标替换为了`to be or not to be`，这样的运行时间会少一点。在主函数中我们经历了几代人，对于每代人我们都检查最佳DNA是否与当代人匹配。如果不匹配，就继续繁殖下一代，如果匹配就退出。

找到每一代中的最优个体

```go
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
```

我们的所有代码都已经写好了，运行一下看看你用了多长时间？

![](C:\Users\lixin yang\Desktop\na.png)

事实上，如果没有濒临灭绝机制的话，适应度`fitness`理论上是不会减少的，但我们的灭绝机制导致运行时间过长。

代码我放在这儿，有需要的朋友可以自行查看

https://github.com/LixvYang/GAS/blob/main/main.go