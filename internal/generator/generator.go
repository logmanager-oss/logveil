package generator

import (
	"github.com/go-faker/faker/v4"
)

type Generator struct{}

func (g *Generator) GenerateRandomIPv4() string {
	return faker.IPv4()
}

func (g *Generator) GenerateRandomIPv6() string {
	return faker.IPv6()
}

func (g *Generator) GenerateRandomMac() string {
	return faker.MacAddress()
}

func (g *Generator) GenerateRandomEmail() string {
	return faker.Email()
}

func (g *Generator) GenerateRandomUrl() string {
	return faker.URL()
}
