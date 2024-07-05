package utils

import "math/rand"


func Map[T, I any](s []T, fn func(T) I ) []I{
	
	var res []I = make([]I,len(s))
	for i, item := range s{
		res[i] = fn(item)
	}
	return res
}

func RandElement[T any](s []T) T{
	random_index := rand.Intn(len(s))
	return s[random_index]
}