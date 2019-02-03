package analyze

import "fmt"

// A the way to pass data to functions
type A func([]byte) (n int, err error)

// Print contents of process
func Print(p []byte) (n int, err error) {

	fmt.Printf("%s\n", p)
	return 0, nil
}
