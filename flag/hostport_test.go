package flag

import (
	"testing"

	"github.com/turbinelabs/test/assert"
)

func TestHostPortInterfaces(t *testing.T) {
	hp := &HostPort{}

	assert.DeepEqual(t, hp.Get(), hp)
	assert.Equal(t, hp.String(), ":0")
	assert.Equal(t, hp.Addr(), ":0")
	h, p := hp.ParsedHostPort()
	assert.Equal(t, h, "")
	assert.Equal(t, p, 0)

	assert.Nil(t, hp.Set(":123"))
	assert.Equal(t, hp.String(), ":123")
	assert.Equal(t, hp.Addr(), ":123")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "")
	assert.Equal(t, p, 123)

	assert.Nil(t, hp.Set("localhost:456"))
	assert.Equal(t, hp.String(), "localhost:456")
	assert.Equal(t, hp.Addr(), "localhost:456")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "localhost")
	assert.Equal(t, p, 456)

	assert.Nil(t, hp.Set("localhost:http"))
	assert.Equal(t, hp.String(), "localhost:http")
	assert.Equal(t, hp.Addr(), "localhost:http")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "localhost")
	assert.Equal(t, p, 80)

	assert.Nil(t, hp.Set("[::1]:http"))
	assert.Equal(t, hp.String(), "[::1]:http")
	assert.Equal(t, hp.Addr(), "[::1]:http")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "::1")
	assert.Equal(t, p, 80)

	assert.Nil(t, hp.Set("fred%zone:http"))
	assert.Equal(t, hp.String(), "fred%zone:http")
	assert.Equal(t, hp.Addr(), "fred%zone:http")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "fred%zone")
	assert.Equal(t, p, 80)

	assert.Nil(t, hp.Set("example.com:http"))
	assert.Equal(t, hp.String(), "example.com:http")
	assert.Equal(t, hp.Addr(), "example.com:http")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "example.com")
	assert.Equal(t, p, 80)

	hp = nil
	assert.Equal(t, hp.String(), ":0")
	assert.Equal(t, hp.Addr(), ":0")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "")
	assert.Equal(t, p, 0)

	hp = &HostPort{host: "host"}
	assert.Equal(t, hp.String(), "host:0")
	assert.Equal(t, hp.Addr(), "host:0")
	h, p = hp.ParsedHostPort()
	assert.Equal(t, h, "host")
	assert.Equal(t, p, 0)

	assert.ErrorContains(t, hp.Set("localhost"), "missing port")
	assert.ErrorContains(t, hp.Set("localhost:99999"), "invalid port")

	// Supremely generic error checking since this varies by OS.
	assert.ErrorContains(t, hp.Set("localhost:portyportport"), "tcp/portyportport")
}

func TestNewHostPort(t *testing.T) {
	hp := NewHostPort("valid:99")
	assert.Equal(t, hp.Addr(), "valid:99")

	hp = NewHostPort("nope")
	assert.Equal(t, hp.Addr(), ":0")
}

func TestNewHostPortWithDefaultPort(t *testing.T) {
	{
		callbacks := 0
		hp := NewHostPortWithDefaultPort("example.com", func() int {
			callbacks++
			return 999
		})

		assert.Equal(t, callbacks, 0)
		assert.Equal(t, hp.Addr(), "example.com:999")
		assert.Equal(t, callbacks, 1)
		assert.Equal(t, hp.Addr(), "example.com:999")
		assert.Equal(t, callbacks, 2)
	}

	{
		callbacks := 0
		hp := NewHostPortWithDefaultPort("example.com", func() int {
			callbacks++
			return 999
		})

		assert.Nil(t, hp.Set("localhost"))
		assert.Equal(t, callbacks, 0)

		assert.Equal(t, hp.Addr(), "localhost:999")
		assert.Equal(t, callbacks, 1)
	}

	{
		hp := NewHostPortWithDefaultPort("example.com", func() int {
			t.Error("unexpected call")
			return 999
		})

		assert.Nil(t, hp.Set("localhost:123"))
		assert.Equal(t, hp.Addr(), "localhost:123")
	}

	{
		hp := NewHostPortWithDefaultPort("[::1]", func() int {
			t.Error("unexpected call")
			return 999
		})

		assert.Nil(t, hp.Set("[fe80::1]:123"))
		assert.Equal(t, hp.Addr(), "[fe80::1]:123")
	}

	{
		hp := NewHostPortWithDefaultPort("localhost", func() int {
			t.Error("unexpected call")
			return 999
		})

		assert.ErrorContains(t, hp.Set("[::1:123"), "missing ']' in address")
	}
}

func TestFlagSetHostPort(t *testing.T) {
	fs := NewTestFlagSet()

	hp := fs.HostPort("addr", NewHostPort("localhost:123"), "usage")

	assert.Nil(t, fs.Parse([]string{"--addr", "example.com:http"}))
	assert.Equal(t, hp.Addr(), "example.com:http")
}

func TestFlagSetHostPortVar(t *testing.T) {
	fs := NewTestFlagSet()

	hp := NewHostPort("localhost:http")
	fs.HostPortVar(&hp, "addr", hp, "usage")

	assert.Nil(t, fs.Parse([]string{"--addr", "example.com:80"}))
	assert.Equal(t, hp.Addr(), "example.com:80")
}
