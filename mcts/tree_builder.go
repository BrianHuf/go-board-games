package mcts

// FixedIterations ...
func FixedIterations(iterations int) (builder TreeBuilderType) {
	builder = func(root Node, config Config) {
		BuildFixedIterations(root, iterations, config)
	}
	return
}

// BuildFixedIterations ...
func BuildFixedIterations(root Node, iterations int, config Config) {
	for i := 0; i < iterations; i++ {
		DoRound(root, config)
	}
}
