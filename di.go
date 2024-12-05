package di

import "reflect"

type dependencyContainer struct {
	bindings   map[reflect.Type]reflect.Type
	dependencies map[reflect.Type]reflect.Value
}

func NewDependencyContainer() *dependencyContainer {
	return &dependencyContainer{
		bindings: make(map[reflect.Type]reflect.Type),
		dependencies: make(map[reflect.Type]reflect.Value),
	}
}

func Bind[T interface{}, K interface{}](container *dependencyContainer) {
	_interfaceType := (*T)(nil)
	_implType := (*K)(nil)

	interfaceType := reflect.TypeOf(_interfaceType).Elem()
	implType := reflect.TypeOf(_implType)

	if interfaceType.Kind() != reflect.Interface {
		panicFirstArgumentNotInterface()
	}
	if !implType.Implements(interfaceType) {
		panicTypeDoesntImplementInterface(implType.Name(), interfaceType.Name())
	}

	container.bindings[interfaceType] = implType
}

func InitHandler[T any](container *dependencyContainer) *T {
	handlerObjRaw := (*T)(nil)
	handlerType := reflect.TypeOf(handlerObjRaw).Elem()

	if handlerType.Kind() == reflect.Interface {
		panicHandlerShouldntBeAnInterface()
	}

	handlerValue := reflect.New(handlerType)
	container.provide(handlerValue.Interface())

	castedHandler, ok := handlerValue.Elem().Interface().(T)
	if !ok {
		panicTypeCastFailed()
	}

	return &castedHandler
}

func (c *dependencyContainer) provide(object any) {
	objectType := reflect.TypeOf(object)
	construct, exists := objectType.MethodByName("Construct")

	if !exists {
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
			if injectionType, exists = c.bindings[injectionType]; !exists {
				panicBindingNotFound(injectionType.Name())
			}
		}

		if injectionValue, exists := c.dependencies[injectionType]; exists {
			constructArgs = append(constructArgs, injectionValue)
			continue
		} 
			
		injectionValue := reflect.New(injectionType.Elem())
		c.provide(injectionValue.Interface())
		constructArgs = append(constructArgs, injectionValue)

		c.dependencies[injectionType] = injectionValue
	}

	construct.Func.Call(
		append([]reflect.Value{reflect.ValueOf(object)}, constructArgs...),
	)
}
