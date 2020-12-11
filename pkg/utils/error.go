package utils

// SecurePanic trigger panic when err is not nil
func SecurePanic(err error)  {
	if err != nil { panic(err)}
}
