package filters

import (
	. "spexs"
	"errors"
)

func trueFilter(p *Pattern, ref *Reference) bool {
	return true
}

func Compose(conf map[string]Conf) (Func, error) {
	filters := make([]Func, 0)

	for name, args := range conf {
		filter, valid := Get(name)
		if !valid {
			return nil, errors.New("No filter named: " + name)
		}
		
		f, err := filter.Create(args)
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
	return func(p *Pattern, ref *Reference) bool {
		for _, f := range filters {
			if !f(p, ref) {
				return false
			}
		}
		return true
	}, nil
}