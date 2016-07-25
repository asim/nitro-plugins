package nats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/micro/go-micro/registry"
)

func TestRegister(t *testing.T) {
	service := registry.Service{Name: "test"}
	require.NoError(t, e.registryOne.Register(&service))
	defer e.registryOne.Deregister(&service)

	services, err := e.registryOne.ListServices()
	require.NoError(t, err)
	assert.Equal(t, 3, len(services))

	services, err = e.registryTwo.ListServices()
	require.NoError(t, err)
	assert.Equal(t, 3, len(services))
}

func TestDeregister(t *testing.T) {
	t.Skip("not properly implemented")

	service := registry.Service{Name: "test"}

	require.NoError(t, e.registryOne.Register(&service))
	require.NoError(t, e.registryOne.Deregister(&service))

	services, err := e.registryOne.ListServices()
	require.NoError(t, err)
	assert.Equal(t, 0, len(services))

	services, err = e.registryTwo.ListServices()
	require.NoError(t, err)
	assert.Equal(t, 0, len(services))
}

func TestGetService(t *testing.T) {
	services, err := e.registryTwo.GetService("one")
	require.NoError(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, "one", services[0].Name)
	assert.Equal(t, 1, len(services[0].Nodes))
}

func TestGetServiceWithNoNodes(t *testing.T) {
	services, err := e.registryOne.GetService("missing")
	require.NoError(t, err)
	assert.Equal(t, 0, len(services))
}

func TestGetServiceFromMultipleNodes(t *testing.T) {
	services, err := e.registryOne.GetService("two")
	require.NoError(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, "two", services[0].Name)
	assert.Equal(t, 2, len(services[0].Nodes))
}

func BenchmarkGetService(b *testing.B) {
	for n := 0; n < b.N; n++ {
		services, err := e.registryTwo.GetService("one")
		require.NoError(b, err)
		assert.Equal(b, 1, len(services))
		assert.Equal(b, "one", services[0].Name)
	}
}

func BenchmarkGetServiceWithNoNodes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		services, err := e.registryOne.GetService("missing")
		require.NoError(b, err)
		assert.Equal(b, 0, len(services))
	}
}

func BenchmarkGetServiceFromMultipleNodes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		services, err := e.registryTwo.GetService("two")
		require.NoError(b, err)
		assert.Equal(b, 1, len(services))
		assert.Equal(b, "two", services[0].Name)
		assert.Equal(b, 2, len(services[0].Nodes))
	}
}
