package generator

import (
	"github.com/go-faker/faker/v4"
)

type Generator struct{}

func (g *Generator) GenerateRandomIPv4() string {
	return faker.IPv4()
}
