package filters

import (
	"errors"
	. "spexs"
)

func trueFilter(p *Query) bool {
	return true
}

func Compose(conf map[string]Conf, setup Setup) (Func, error) {
	filters := make([]Func, 0)

	for name, args := range conf {
		filter, valid := Get(name)
		if !valid {
			return nil, errors.New("No filter named: " + name)
		}

		f, err := filter.Create(args, setup)
		if err != nil {
			return nil, err
		}

		filters = append(filters, f)
	}

	if len(filters) == 0 {
		return trueFilter, nil
	} else if len(filters) == 1 {
		return filters[0], nil
	}

	// create a composite filter
	return func(p *Query) bool {
		for _, f := range filters {
			if !f(p) {
				return false
			}
		}
		return true
	}, nil
}
