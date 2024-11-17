//go:generate bash -c "mockgen -source=$(basename ${GOFILE} .go).go -package=$(go list -f '{{.Name}}') -destination=$(basename ${GOFILE} .go)_mock_test.go"
package dog

type Dog struct {
	ID    string
	Name  string
	Breed string
}

type DogDatabase interface {
	Insert(Dog) error
	SelectAll() ([]Dog, error)
}
