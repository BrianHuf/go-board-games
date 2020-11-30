package mcts

// FixedIterations ...
func FixedIterations(iterations int) (builder TreeBuilderType) {
	builder = func(root Node, config Config) {
		buildFixedIterations(root, iterations, config)
	}
	return
}

// FixedIterationsParallel ...
func FixedIterationsParallel(iterations int, parallel int) (builder TreeBuilderType) {
	builder = func(root Node, config Config) {
		buildFixedIterationsParallel(root, iterations, parallel, config)
	}
	return
}

func buildFixedIterations(root Node, iterations int, config Config) {
	for i := 0; i < iterations; i++ {
		selected := Selector(root, config)
		new := Expand(selected, config)
		score, _ := Playout(new.GetMove(), config)
		Propagate(new, score, config)
	}
}

type updateNode struct {
	node  Node
	score float32
}

func buildFixedIterationsParallel(root Node, iterations int, parallel int, config Config) {
	updateChannel := make(chan updateNode, parallel)
	playoutChannel := make(chan Node, parallel)

	// parallel Playouts
	for i := 0; i < parallel; i++ {
		go func() {
			for playoutNode := range playoutChannel {
				score, _ := Playout(playoutNode.GetMove(), config)
				updateChannel <-updateNode{playoutNode, score}
			}
		}()
	}

	// run 
	for i := 0; i< parallel; i++ {
		updateChannel <-updateNode{}
	}
	
	count := 0
	for updateNode := range updateChannel {
		if updateNode.node != nil {
			Propagate(updateNode.node, updateNode.score, config)
			count++
		}

		if (count == iterations) {
			break
		}

		selected := Selector(root, config)
		newNode := Expand(selected, config)
		playoutChannel <-newNode
	}
}
