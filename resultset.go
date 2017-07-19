package redirection

// ResultSet is a set of results.
type ResultSet []Result

func (set ResultSet) Len() int           { return len(set) }
func (set ResultSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set ResultSet) Less(i, j int) bool { return set[i].Weight > set[j].Weight }
