package main

import (
	"strings"
	"sync"
)

var (
	roleMaxId = 10
)

type RoleTrie []*Trie

type Trie struct {
	rw       sync.RWMutex
	RoleId   int      `json:"role_id"`
	RoleName string   `json:"role_name"`
	Node     TrieNode `json:"node"`
}

type TrieNode map[string]TrieNodeInfo

type TrieNodeInfo struct {
	Word     string   `json:"word"`
	Children TrieNode `json:"children,omitempty"`
	Method   string   `json:"method,omitempty"`
	IsLast   bool     `json:"is_last,omitempty"`
	Name     string   `json:"name,omitempty"`
}

func NewRoleTrie() RoleTrie {
	return make([]*Trie, roleMaxId)
}

type RoleInfo struct {
	RoleId   int       `json:"role_id"`
	RoleName string    `json:"role_name"`
	Apis     []ApiInfo `json:"apis"`
}

type ApiInfo struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

func (r RoleTrie) Generate(roleInfo RoleInfo) {
	trie := r[roleInfo.RoleId]
	if trie == nil {
		trie = &Trie{
			RoleId:   roleInfo.RoleId,
			RoleName: roleInfo.RoleName,
		}
	}

	trie.rw.Lock() // 这里可能加锁不成功，暂时忽略这个问题
	defer trie.rw.Unlock()

	trie.Node = make(map[string]TrieNodeInfo)
	for _, api := range roleInfo.Apis {
		arr := strings.Split(api.Url, "/")
		if len(arr) == 0 {
			continue
		}

		node := NewNode(trie, api, arr, 0)
		GenerateTrieNodeInfo(node, api, arr, 1)
		trie.Node[arr[0]] = node
	}
	r[roleInfo.RoleId] = trie
}

func (r RoleTrie) Search(roleId int, url, method string) bool {
	arr := strings.Split(url, "/")
	roleInfo := r[roleId]
	roleInfo.rw.RLock()
	defer roleInfo.rw.RUnlock()
	if len(arr) == 0 {
		return false
	}
	node, ok := roleInfo.Node[arr[0]]
	if !ok {
		return false
	}
	if len(arr) == 1 {
		if node.IsLast == true && node.Method == method {
			return true
		} else {
			return false
		}
	}

	for i := 1; i < len(arr); i++ {
		tmp, ok := node.Children[arr[i]]
		if !ok {
			if tmp, ok = node.Children[":str"]; !ok {
				return false
			}
		}
		if ok {
			if len(arr) == i+1 {
				if tmp.IsLast == true && tmp.Method == method {
					return true
				} else {
					return false
				}
			} else {
				node = tmp
			}
		}
	}
	return false
}

func GenerateTrieNodeInfo(data TrieNodeInfo, api ApiInfo, arr []string, index int) {
	dataNode := NewNodeInfo(data, api, arr, index)
	if len(arr)-1 != index {
		GenerateTrieNodeInfo(dataNode, api, arr, index+1)
	}
	data.Children[arr[index]] = dataNode
	return
}

func NewNodeInfo(data TrieNodeInfo, api ApiInfo, arr []string, index int) TrieNodeInfo {
	dataNode, ok := data.Children[arr[index]]
	if !ok {
		dataNode = TrieNodeInfo{
			Word:     arr[index],
			Children: make(map[string]TrieNodeInfo),
		}
	}
	if len(arr)-1 == index {
		dataNode.Children = nil
		dataNode.Method = api.Method
		dataNode.IsLast = true
		dataNode.Name = api.Name
	}
	return dataNode
}

func NewNode(trie *Trie, api ApiInfo, arr []string, index int) TrieNodeInfo {
	node, ok := trie.Node[arr[0]]
	if !ok {
		node = TrieNodeInfo{
			Word:     arr[0],
			Children: make(map[string]TrieNodeInfo),
		}
	}
	if len(arr)-1 == index {
		node.Children = nil
		node.Method = api.Method
		node.IsLast = true
		node.Name = api.Name
	}
	return node
}
