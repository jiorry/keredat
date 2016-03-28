package lib

import "sort"

type SortBy func(p1, p2 *ExcelRow) bool

func (by SortBy) Sort(rows []ExcelRow) {
	ps := &Sorter{
		rows: rows,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// Sorter joins a By function and a slice of Planets to be sorted.
type Sorter struct {
	rows []ExcelRow
	by   func(p1, p2 *ExcelRow) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *Sorter) Len() int {
	return len(s.rows)
}

// Swap is part of sort.Interface.
func (s *Sorter) Swap(i, j int) {
	s.rows[i], s.rows[j] = s.rows[j], s.rows[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *Sorter) Less(i, j int) bool {
	return s.by(&s.rows[i], &s.rows[j])
}
