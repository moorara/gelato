package visitor

import "github.com/moorara/gelato/internal/log"

type visitor struct {
	depth  int
	logger *log.ColorfulLogger
}

func (v visitor) Next() visitor {
	return visitor{
		depth:  v.depth + 1,
		logger: v.logger,
	}
}
