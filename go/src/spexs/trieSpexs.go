package spexs

type TrieNodes chan *TrieNode

func MakeTrieNodes() TrieNodes {
	return make(TrieNodes, patternsBufferSize)
}

type TriePooler interface {
	Take() (*TrieNode, bool)
	Put(*TrieNode)
	Len() int
}

type TrieFilterFunc func(p *TrieNode) bool
type TrieExtenderFunc func(p *TrieNode, ref *UnicodeReference) TrieNodes

func RunTrie(ref *UnicodeReference, input TriePooler, results TriePooler,
	extender TrieExtenderFunc, acceptable TrieFilterFunc) {
	p, valid := input.Take()
	for valid {
		pats := extender(p, ref)
		for ep := range pats {
			if acceptable(ep) {
				input.Put(ep)
				results.Put(ep)
			}
		}
		p, valid = input.Take()
	}
}

func RunTrieParallel(ref *UnicodeReference, input TriePooler, results TriePooler,
	extender TrieExtenderFunc, acceptable TrieFilterFunc, num_threads int) {

	start := make(chan int, 1000)
	stop := make(chan int, 1000)
	
	for i := 0; i < num_threads; i++ {
		go func(){
			start <- 1
			defer func() { stop <- 1 }()

			for {
				p, valid := input.Take()
				if !valid {
					return
				}

				pats := extender(p, ref)
				for ep := range pats {
					if acceptable(ep) {
						input.Put(ep)
						results.Put(ep)
					}
				}
			}
		}()
	}

	for i := 0; i < num_threads; i++ {
		<-stop
	}
}
