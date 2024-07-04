package utils


func Map[T, I any](s []T, fn func(T) I ) []I{
	
	var res []I = make([]I,len(s))
	for i, item := range s{
		res[i] = fn(item)
	}
	return res
}