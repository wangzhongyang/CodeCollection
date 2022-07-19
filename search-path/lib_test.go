package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

var (
	roleTrieTest RoleTrie
	roleRegTest  RoleReg
	url1         = "github.com/go-redis/nihao/redis/v8"
	url2         = "github.com/go-redis/wohao/redis/v8"
	apiStr       = `{"role_id":5,"role_name":"超级管理员","apis":[{"name":"a","url":"github.com/go-redis/:str/redis/v8","method":"GET"},{"name":"b","url":"golang.org/x/net/:str","method":"GET"},{"name":"c","url":"github.com/cespare/xxhash/v2","method":"POST"},{"name":"d","url":"github.com/dgryski/go-rendezvous/:str","method":"POST"},{"name":"e","url":"github.com/:str/go-redis/redis/v8","method":"GET"},{"name":"f","url":"golang.org/:str/x/net","method":"GET"},{"name":"g","url":"gopkg.in/yaml.v2/asds/:str","method":"GET"},{"name":"h","url":"golang.org/x/text/:str","method":"GET"},{"name":"i","url":"github.com/davecgh/go-spew/:str","method":"GET"},{"name":"j","url":"github.com/envoyproxy/go-control-plane/:str","method":"GET"},{"name":"k","url":"github.com/go-playground/validator/v10/:str","method":"GET"},{"name":"l","url":"github.com/grpc/ecosystem/grpc/gateway/:str","method":"GET"}]}`
)

func TestMain(m *testing.M) {
	roleInfo := new(RoleInfo)
	_ = json.Unmarshal([]byte(apiStr), roleInfo)

	roleTrieTest = NewRoleTrie()
	roleTrieTest.Generate(*roleInfo)

	roleRegTest = NewRegexp()
	roleRegTest.GenerateReg(*roleInfo)

	m.Run()
}

func Benchmark_TrieSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if roleTrieTest.Search(5, url1, http.MethodGet) != true {
			b.Fatal("Benchmark_TrieSearch url1 has failed")
		}
		if roleTrieTest.Search(5, url2, http.MethodPost) != false {
			b.Fatal("Benchmark_TrieSearch url2 has failed")
		}
	}
}

func Benchmark_RegSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if roleRegTest.Search(5, url1, http.MethodGet) != true {
			b.Fatal("Benchmark_RegSearch url1 has failed")
		}
		if roleRegTest.Search(5, url2, http.MethodPost) != false {
			b.Fatal("Benchmark_RegSearch url2 has failed")
		}
	}
}
