package main

import (
	"di/handlers"
	"di/repos"
	"fmt"
	"reflect"
)

func main() {
	di := NewDependencyContainer()

	Bind[repos.Repo, repos.RepoImpl](di)

	postHandler := InitHandler[handlers.PostHandler](di).Handle

	postHandler()
}

type DependencyContainer struct {
	bindings   map[reflect.Type]reflect.Type
	Dependencies map[reflect.Type]reflect.Value
}

func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		bindings: make(map[reflect.Type]reflect.Type),
		Dependencies: make(map[reflect.Type]reflect.Value),
	}
}

func (c *DependencyContainer) GetBinding(interfaceType reflect.Type) reflect.Type {
	implType, exists := c.bindings[interfaceType]
	if !exists {
		panic(fmt.Sprintf("Для интерфейса %s не зарегистрирована реализация", interfaceType))
	}
	return implType
}

func Bind[T interface{}, K interface{}](container *DependencyContainer) {
	_iType := (*T)(nil)
	_implType := (*K)(nil)

	iType := reflect.TypeOf(_iType).Elem()
	implType := reflect.TypeOf(_implType)

	if iType.Kind() != reflect.Interface {
		panic("Первый аргумент должен быть интерфейсом")
	}
	if !implType.Implements(iType) {
		panic(fmt.Sprintf("Тип %s не реализует интерфейс %s", implType, iType))
	}

	container.bindings[iType] = implType
}

func InitHandler[T any](container *DependencyContainer) *T {
	handlerObjRaw := (*T)(nil)
	handlerType := reflect.TypeOf(handlerObjRaw).Elem()

	if handlerType.Kind() == reflect.Interface {
		panic("Handler не может быть интерфейсом")
	}

	handlerValue := reflect.New(handlerType)
	container.Provide(handlerValue.Interface())

	castedHandler, ok := handlerValue.Elem().Interface().(T)
	if !ok {
		panic("Не удалось произвести приведение типов при инициализации Handler")
	}

	return &castedHandler
}

func (c *DependencyContainer) Provide(object any) {
	objectType := reflect.TypeOf(object)
	construct, exists := objectType.MethodByName("Construct")

	if !exists {
		fmt.Println("no constructor")
		return
	}

	constructArgs := make([]reflect.Value, 0, construct.Type.NumIn()-1)

	for i := 1; i < construct.Type.NumIn(); i++ {
		injectionType := construct.Type.In(i)

		// инъектить можно только либо указатель, либо интерфейс
		if injectionType.Kind() != reflect.Ptr && injectionType.Kind() != reflect.Interface {
			continue
		}

		if injectionType.Kind() == reflect.Interface {
			// ищем в контейнере имплементацию интерфейса
			// подставляем тип имплементации
			injectionType = c.GetBinding(injectionType)
		}

		if injectionValue, exists := c.Dependencies[injectionType]; exists {
			constructArgs = append(constructArgs, injectionValue)
			fmt.Printf(
				"injectionValue of type: %s got from Dependencies while constructing: %s\n",
				injectionType.Elem().Name(),
				objectType.Elem().Name(),
			)
			continue
			} 
			
		injectionValue := reflect.New(injectionType.Elem())
		c.Provide(injectionValue.Interface())
		constructArgs = append(constructArgs, injectionValue)

		c.Dependencies[injectionType] = injectionValue
	}

	construct.Func.Call(
		append([]reflect.Value{reflect.ValueOf(object)}, constructArgs...),
	)
	fmt.Printf("construct called in: %s\n", objectType.Elem().Name())
}
