package dedupe

// Interface, if implemented, enables a struct to be used to remove entry
// duplications in a slice. It is assumed that the slice the Interface allows
// interaction with has all like objects grouped already. This need not be
// sorted by that trivially fulfills the requirement.
type Interface interface {
	Len() int
	Equal(i, j int) bool
	Remove(i int)
}

// Dedupe removes duplicate entries from within a slice represented by the
// provided Interface.
func Dedupe(tgt Interface) {
	exitLen := tgt.Len()
	for i := 0; i < exitLen; exitLen = tgt.Len() {
		until := i
		for until < exitLen && tgt.Equal(i, until) {
			until++
		}

		if i != until {
			// found some dupes
			for remove := (until - 1); remove > i; remove-- {
				tgt.Remove(remove)
			}
		}
		i++
	}
}

type ints struct {
	is []int
}

func (a *ints) Len() int            { return len(a.is) }
func (a *ints) Equal(i, j int) bool { return a.is[i] == a.is[j] }
func (a *ints) Remove(i int)        { a.is = append(a.is[0:i], a.is[i+1:]...) }

// Ints returns an array of ints that has the duplicated entries removed. It
// is assumed that the input []int has all equal values grouped.
func Ints(i []int) []int {
	is := &ints{i}
	Dedupe(is)
	return is.is
}

type strings struct {
	ss []string
}

func (a *strings) Len() int            { return len(a.ss) }
func (a *strings) Equal(i, j int) bool { return a.ss[i] == a.ss[j] }
func (a *strings) Remove(i int)        { a.ss = append(a.ss[0:i], a.ss[i+1:]...) }

// Strings returns an array of strings that has the duplicated entries removed.
// It is assumed that the input []string has all equal values grouped.
func Strings(s []string) []string {
	ss := &strings{s}
	Dedupe(ss)
	return ss.ss
}
