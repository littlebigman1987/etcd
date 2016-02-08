// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
)

// GlobalFlags are flags that defined globally
// and are inherited to all sub-commands.
type GlobalFlags struct {
	Endpoints string
	TLS       transport.TLSInfo
}

func mustClient(cmd *cobra.Command) *clientv3.Client {
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		ExitWithError(ExitError, err)
	}

	// set tls if any one tls option set
	var cfgtls *transport.TLSInfo
	tls := transport.TLSInfo{}
	var file string
	if file, err = cmd.Flags().GetString("cert"); err == nil && file != "" {
		tls.CertFile = file
		cfgtls = &tls
	} else if cmd.Flags().Changed("cert") {
		ExitWithError(ExitBadArgs, errors.New("empty string is passed to --cert option"))
	}

	if file, err = cmd.Flags().GetString("key"); err == nil && file != "" {
		tls.KeyFile = file
		cfgtls = &tls
	} else if cmd.Flags().Changed("key") {
		ExitWithError(ExitBadArgs, errors.New("empty string is passed to --key option"))
	}

	if file, err = cmd.Flags().GetString("cacert"); err == nil && file != "" {
		tls.CAFile = file
		cfgtls = &tls
	} else if cmd.Flags().Changed("cacert") {
		ExitWithError(ExitBadArgs, errors.New("empty string is passed to --cacert option"))
	}

	cfg := clientv3.Config{
		Endpoints:   []string{endpoint},
		TLS:         cfgtls,
		DialTimeout: 20 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		ExitWithError(ExitBadConnection, err)
	}
	return client
}

func argOrStdin(args []string, stdin io.Reader, i int) ([]byte, error) {
	if i < len(args) {
		return []byte(args[i]), nil
	}
	bytes, err := ioutil.ReadAll(stdin)
	if string(bytes) == "" || err != nil {
		return nil, errors.New("no available argument and stdin")
	}
	return bytes, nil
}
