package main

import (
	"regexp"
	"strings"
)

type RoleReg []RoleRegItem

type RoleRegItem struct {
	Method  string        `json:"method"`
	Name    string        `json:"name"`
	RoleId  int           `json:"role_id"`
	RegPath []RegPathInfo `json:"reg_path"`
}

type RegPathInfo struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func NewRegexp() RoleReg {
	return make([]RoleRegItem, roleMaxId)
}

func (r RoleReg) GenerateReg(roleInfo RoleInfo) {
	r[roleInfo.RoleId] = RoleRegItem{
		Name:   roleInfo.RoleName,
		RoleId: roleInfo.RoleId,
	}
	for _, val := range roleInfo.Apis {
		r[roleInfo.RoleId].RegPath = append(r[roleInfo.RoleId].RegPath, RegPathInfo{
			Path:   strings.ReplaceAll(val.Url, ":str", "([1-9]\\d*|[a-zA-Z]*)"),
			Name:   val.Name,
			Method: val.Method,
		})
	}

	return
}

func (r RoleReg) Search(roleId int, url, method string) bool {
	roleItem := r[roleId]
	if roleItem.Name == "" {
		return false
	}
	for _, v := range roleItem.RegPath {
		if regexp.MustCompile(v.Path).MatchString(url) && method == v.Method {
			return true
		}
	}
	return false
}
