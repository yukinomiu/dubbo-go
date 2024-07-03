/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 *
 * Copyright 2021 gRPC authors.
 *
 */

// Package unix implements a resolver for unix targets.
package unix

import (
	"fmt"
)

import (
	"google.golang.org/grpc/resolver"
)

import (
	"dubbo.apache.org/dubbo-go/v3/xds/utils/transport/networktype"
)

const unixScheme = "unix"
const unixAbstractScheme = "unix-abstract"

type builder struct {
	scheme string
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	// dubbox fix
	if authority := target.URL.Host; authority != "" {
		return nil, fmt.Errorf("invalid (non-empty) authority: %v", authority)
	}

	// gRPC was parsing the dial target manually before PR #4817, and we
	// switched to using url.Parse() in that PR. To avoid breaking existing
	// resolver implementations we ended up stripping the leading "/" from the
	// endpoint. This obviously does not work for the "unix" scheme. Hence we
	// end up using the parsed URL instead.
	endpoint := target.URL.Path
	if endpoint == "" {
		endpoint = target.URL.Opaque
	}
	addr := resolver.Address{Addr: endpoint}
	if b.scheme == unixAbstractScheme {
		// prepend "\x00" to address for unix-abstract
		addr.Addr = "\x00" + addr.Addr
	}
	cc.UpdateState(resolver.State{Addresses: []resolver.Address{networktype.Set(addr, "unix")}})
	return &nopResolver{}, nil
}

func (b *builder) Scheme() string {
	return b.scheme
}

type nopResolver struct {
}

func (*nopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (*nopResolver) Close() {}

func init() {
	resolver.Register(&builder{scheme: unixScheme})
	resolver.Register(&builder{scheme: unixAbstractScheme})
}
