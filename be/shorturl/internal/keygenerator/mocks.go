package keygenerator

import "github.com/stretchr/testify/mock"

type (
	MockGenerator struct {
		mock.Mock
	}
)

// Generate implements KeyGenerator.
func (m *MockGenerator) Generate() string {
	args := m.Called()
	return args.Get(0).(string)
}

var _ (KeyGenerator) = (*MockGenerator)(nil)
