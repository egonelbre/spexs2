package spexs

type FifoPool struct {
  patterns Patterns
}

func (p *FifoPool) Take() (Pattern, bool) {
  return <- p.patterns
}

func (p *FifoPool) Put( pat Pattern ) {
  p.patterns <- pat
}