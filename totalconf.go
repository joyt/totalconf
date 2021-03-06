package totalconf

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/rakyll/globalconf"
)

type Options struct {
	globalconf.Options
	EtcdClient *etcd.Client
}

var (
	flags        = map[string]*flag.FlagSet{}
	mu           sync.Mutex
	parsed       bool
	onParseFuncs []func()
)

func OnParsed(do func()) {
	mu.Lock()
	defer mu.Unlock()
	if parsed {
		go do()
	} else {
		onParseFuncs = append(onParseFuncs, do)
	}
}

func Parse(opts *Options) error {
	mu.Lock()
	defer mu.Unlock()
	if parsed {
		return nil
	}
	var (
		conf *globalconf.GlobalConf
		err  error
	)
	if opts == nil {
		conf, err = globalconf.New(os.Args[0])
	} else {
		conf, err = globalconf.NewWithOptions(&opts.Options)
	}
	if err != nil {
		return err
	}
	conf.ParseAll()
	set := map[string]bool{}
	flag.Visit(func(f *flag.Flag) {
		set[f.Name] = true
	})
	for name, flagset := range flags {
		if set[name] {
			val := flag.Lookup(name).Value.String()
			flagset.VisitAll(func(f *flag.Flag) {
				f.Value.Set(val)
			})
		} else if opts.EtcdClient != nil {
			if resp, err := opts.EtcdClient.Get(name, false, false); err == nil && resp.Node != nil {
				flagset.VisitAll(func(f *flag.Flag) {
					f.Value.Set(resp.Node.Value)
				})
			}
		}
	}
	parsed = true
	for _, f := range onParseFuncs {
		go f()
	}
	return nil
}

func Parsed() bool {
	return parsed
}

func String(name string, value string, usage string) *string {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.String(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].String(scope, value, usage)
}

func Bool(name string, value bool, usage string) *bool {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.Bool(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].Bool(scope, value, usage)
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	name = strings.ToLower(name)
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.Duration(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].Duration(scope, value, usage)
}

func Float64(name string, value float64, usage string) *float64 {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.Float64(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].Float64(scope, value, usage)
}

func Int(name string, value int, usage string) *int {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.Int(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].Int(scope, value, usage)
}

func Int64(name string, value int64, usage string) *int64 {
	mu.Lock()
	defer mu.Unlock()
	name = strings.ToLower(name)
	var scope string
	if _, file, line, ok := runtime.Caller(1); ok {
		scope = fmt.Sprintf("%s:%d", file, line)
	}
	if flags[name] == nil {
		flag.Int64(name, value, usage)
		flags[name] = flag.NewFlagSet(name, flag.ExitOnError)
	}
	return flags[name].Int64(scope, value, usage)
}
