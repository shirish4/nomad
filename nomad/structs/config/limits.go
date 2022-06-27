package config

import (
	"math"

	"github.com/hashicorp/nomad/helper"
)

const (
	// LimitsNonStreamingConnsPerClient is the number of connections per
	// peer to reserve for non-streaming RPC connections. Since streaming
	// RPCs require their own TCP connection, they have their own limit
	// this amount lower than the overall limit. This reserves a number of
	// connections for Raft and other RPCs.
	//
	// TODO Remove limit once MultiplexV2 is used.
	LimitsNonStreamingConnsPerClient = 20
)

// Limits configures timeout limits similar to Consul's limits configuration
// parameters. Limits is the internal version with the fields parsed.
type Limits struct {
	// HTTPSHandshakeTimeout is the deadline by which HTTPS TLS handshakes
	// must complete.
	//
	// 0 means no timeout.
	HTTPSHandshakeTimeout string `hcl:"https_handshake_timeout"`

	// HTTPMaxConnsPerClient is the maximum number of concurrent HTTP
	// connections from a single IP address. nil/0 means no limit.
	HTTPMaxConnsPerClient *int `hcl:"http_max_conns_per_client"`

	// RPCHandshakeTimeout is the deadline by which RPC handshakes must
	// complete. The RPC handshake includes the first byte read as well as
	// the TLS handshake and subsequent byte read if TLS is enabled.
	//
	// The deadline is reset after the first byte is read so when TLS is
	// enabled RPC connections may take (timeout * 2) to complete.
	//
	// The RPC handshake timeout only applies to servers. 0 means no
	// timeout.
	RPCHandshakeTimeout string `hcl:"rpc_handshake_timeout"`

	// RPCMaxConnsPerClient is the maximum number of concurrent RPC
	// connections from a single IP address. nil/0 means no limit.
	RPCMaxConnsPerClient *int `hcl:"rpc_max_conns_per_client"`

	// RPCDefaultWriteRate is the default maximum write RPC requests
	// per endpoint per user per second. nil/0 means no limit.
	RPCDefaultWriteRate *uint64 `hcl:"rpc_default_write_rate"`

	// RPCDefaultReadRate is the default maximum read RPC requests per
	// endpoint per user per second. nil/0 means no limit.
	RPCDefaultReadRate *uint64 `hcl:"rpc_default_read_rate"`

	// RPCDefaultListRate is the default maximum list RPC requests per
	// endpoint per user per second. nil/0 means no limit.
	RPCDefaultListRate *uint64 `hcl:"rpc_default_list_rate"`

	// These are the RPC limits for individual RPC endpoints
	Namespace *RPCEndpointLimits `hcl:"namespace"`
	Job       *RPCEndpointLimits `hcl:"job"`
	// TODO, etc...
}

// DefaultLimits returns the default limits values. User settings should be
// merged into these defaults.
func DefaultLimits() Limits {
	return Limits{
		HTTPSHandshakeTimeout: "5s",
		HTTPMaxConnsPerClient: helper.IntToPtr(100),
		RPCHandshakeTimeout:   "5s",
		RPCMaxConnsPerClient:  helper.IntToPtr(100),
		RPCDefaultWriteRate:   helper.Uint64ToPtr(math.MaxUint64),
		RPCDefaultReadRate:    helper.Uint64ToPtr(math.MaxUint64),
		RPCDefaultListRate:    helper.Uint64ToPtr(math.MaxUint64),
	}
}

// Merge returns a new Limits where non-empty/nil fields in the argument have
// precedence.
func (l *Limits) Merge(o Limits) Limits {
	m := *l

	if o.HTTPSHandshakeTimeout != "" {
		m.HTTPSHandshakeTimeout = o.HTTPSHandshakeTimeout
	}
	if o.HTTPMaxConnsPerClient != nil {
		m.HTTPMaxConnsPerClient = helper.IntToPtr(*o.HTTPMaxConnsPerClient)
	}
	if o.RPCHandshakeTimeout != "" {
		m.RPCHandshakeTimeout = o.RPCHandshakeTimeout
	}
	if o.RPCMaxConnsPerClient != nil {
		m.RPCMaxConnsPerClient = helper.IntToPtr(*o.RPCMaxConnsPerClient)
	}
	if o.RPCDefaultWriteRate != nil {
		m.RPCDefaultWriteRate = helper.Uint64ToPtr(*o.RPCDefaultWriteRate)
	}
	if o.RPCDefaultReadRate != nil {
		m.RPCDefaultReadRate = helper.Uint64ToPtr(*o.RPCDefaultReadRate)
	}
	if o.RPCDefaultListRate != nil {
		m.RPCDefaultListRate = helper.Uint64ToPtr(*o.RPCDefaultListRate)
	}

	m.Namespace = m.Namespace.Merge(*o.Namespace)

	return m
}

// Copy returns a new deep copy of a Limits struct.
func (l *Limits) Copy() Limits {
	c := *l
	if l.HTTPMaxConnsPerClient != nil {
		c.HTTPMaxConnsPerClient = helper.IntToPtr(*l.HTTPMaxConnsPerClient)
	}
	if l.RPCMaxConnsPerClient != nil {
		c.RPCMaxConnsPerClient = helper.IntToPtr(*l.RPCMaxConnsPerClient)
	}
	if l.RPCDefaultWriteRate != nil {
		c.RPCDefaultWriteRate = helper.Uint64ToPtr(*l.RPCDefaultWriteRate)
	}
	if l.RPCDefaultReadRate != nil {
		c.RPCDefaultReadRate = helper.Uint64ToPtr(*l.RPCDefaultReadRate)
	}
	if l.RPCDefaultListRate != nil {
		c.RPCDefaultListRate = helper.Uint64ToPtr(*l.RPCDefaultListRate)
	}
	c.Namespace = l.Namespace.Copy()

	return c
}

type RPCEndpointLimits struct {
	RPCWriteRate *uint64 `hcl:"rpc_write_rate"`
	RPCReadRate  *uint64 `hcl:"rpc_read_rate"`
	RPCListRate  *uint64 `hcl:"rpc_list_rate"`
}

func (l *RPCEndpointLimits) Merge(o RPCEndpointLimits) *RPCEndpointLimits {
	m := l
	if o.RPCWriteRate != nil {
		m.RPCWriteRate = helper.Uint64ToPtr(*o.RPCWriteRate)
	}
	if o.RPCReadRate != nil {
		m.RPCReadRate = helper.Uint64ToPtr(*o.RPCReadRate)
	}
	if o.RPCListRate != nil {
		m.RPCListRate = helper.Uint64ToPtr(*o.RPCListRate)
	}
	return m
}

func (l *RPCEndpointLimits) Copy() *RPCEndpointLimits {
	c := l
	if l.RPCWriteRate != nil {
		c.RPCWriteRate = helper.Uint64ToPtr(*l.RPCWriteRate)
	}
	if l.RPCReadRate != nil {
		c.RPCReadRate = helper.Uint64ToPtr(*l.RPCReadRate)
	}
	if l.RPCListRate != nil {
		c.RPCListRate = helper.Uint64ToPtr(*l.RPCListRate)
	}
	return c
}
