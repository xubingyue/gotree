// Copyright gotree Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import (
	"errors"
	"reflect"

	"github.com/8treenet/gotree/helper"
)

//服务定位器
type ServiceLocator struct {
	Object
	dict *Dict
}

func (self *ServiceLocator) Gotree() *ServiceLocator {
	self.Object.Gotree(self)
	self.dict = new(Dict).Gotree()
	return self
}

func (self *ServiceLocator) CheckService(com interface{}) bool {
	return self.dict.Check(reflect.TypeOf(com).Elem().Name())
}

//加入服务
func (self *ServiceLocator) AddService(obj interface{}) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		helper.Log().WriteError("AddComponent != reflect.Ptr")
	}
	self.dict.Set(t.Elem().Name(), obj)
}

//移除服务
func (self *ServiceLocator) RemoveService(obj interface{}) {
	t := reflect.TypeOf(obj)
	self.dict.Del(t.Name())
}

//获取服务
func (self *ServiceLocator) Service(obj interface{}) error {
	t := reflect.TypeOf(obj)
	return self.dict.Get(t.Elem().Elem().Name(), obj)
}

//广播定位器内所有实现method的方法
func (self *ServiceLocator) Broadcast(method string, arg interface{}) error {
	list := self.dict.Keys()
	call := false
	for _, v := range list {
		com := self.dict.GetInterface(v)
		if com == nil {
			continue
		}
		value := reflect.ValueOf(com).MethodByName(method)
		if value.Kind() != reflect.Invalid {
			value.Call([]reflect.Value{reflect.ValueOf(arg)})
			call = true
		}
	}
	if !call {
		return errors.New("ServiceLocator-Broadcast Method not found" + method)
	}
	return nil
}
