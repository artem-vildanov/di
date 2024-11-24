package main

import (
	"di/handlers"
	"di/repos"
	"di/services"
	"fmt"
	"reflect"
)

func main() {
	di := NewDependencyContainer()

	Register[repos.RepoInterface, repos.RepoImpl](di)

	Singleton[services.SharedClient](di)

	postHandler := InitHandler[handlers.PostHandler](di).Handle
	userHandler := InitHandler[handlers.UserHandler](di).Handle

	postHandler()
	userHandler()
}

type DependencyContainer struct {
	registry   map[reflect.Type]reflect.Type
	singletons map[reflect.Type]reflect.Value
}

func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		registry: make(map[reflect.Type]reflect.Type),
		singletons: make(map[reflect.Type]reflect.Value),
	}
}

func Register[T interface{}, K interface{}](container *DependencyContainer) {
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

	container.registry[iType] = implType
}

func Singleton[T interface{}](container *DependencyContainer) {
	_implType := (*T)(nil)
	implType := reflect.TypeOf(_implType)

	if implType.Elem().Kind() == reflect.Interface {
		panic("Singleton не может быть интерфейсом")
	}

	if implType.Kind() != reflect.Struct && implType.Kind() != reflect.Ptr {
		panic("Singleton должен быть структурой или указателем на структуру")
	}

	singletonInstance := reflect.New(implType.Elem())
	container.Enrich(singletonInstance.Interface())

	container.singletons[implType] = singletonInstance
}

func (c *DependencyContainer) GetSingleton(searchType reflect.Type) (reflect.Value, bool) {
	value, exists := c.singletons[searchType]
	return value, exists
}

func (c *DependencyContainer) Resolve(interfaceType reflect.Type) reflect.Type {
	implType, exists := c.registry[interfaceType]
	if !exists {
		panic(fmt.Sprintf("Для интерфейса %s не зарегистрирована реализация", interfaceType))
	}
	return implType
}

type Enrichable interface {
	// сработает после того как инициализированы поля структуры
	// можно обратиться к полям и выполнить некоторую логику
	AfterEnriched()
}

type Instantiatable interface {
	// сработает когда был создан начальный экземпляр структуры, и поля еще не инициализированы
	// можно инициализировать поля начальными данными
	// а также можно инициализировать те поля которые обозначены как
	// `di:"enrich"` или `di:"instantiate"`, тем самым
	// преинициализированные поля будут пропущены и не инициализированы повторно
	AfterInstantiated()
}

func InitHandler[T any](container *DependencyContainer) *T {
	handlerObjRaw := (*T)(nil)
	handlerType := reflect.TypeOf(handlerObjRaw).Elem()

	if handlerType.Kind() == reflect.Interface {
		panic("Handler не может быть интерфейсом")
	}

	handlerValue := reflect.New(handlerType)
	container.Enrich(handlerValue.Interface())

	castedHandler, ok := handlerValue.Elem().Interface().(T)
	if !ok {
		panic("Не удалось произвести приведение типов при инициализации Handler")
	}

	return &castedHandler
}

const (
	ignoreTag = "ignore"
	enrichTag = "enrich"
	instantiateTag = "instantiate"
)

// обогащаем инстанциированную структуру
// инициализируя ее поля
// функция вызывается рекурсивно для полей структуры
func (c *DependencyContainer) Enrich(object any) {
	objectValue := reflect.ValueOf(object).Elem()
	objectType := objectValue.Type()

	if method := reflect.ValueOf(object).MethodByName("AfterInstantiated"); 
		method.IsValid() {
		method.Call(nil)
	}

	for i := 0; i < objectType.NumField(); i++ {
		fieldMetadata := objectType.Field(i)
		fieldValue := objectValue.Field(i)
		fieldType := fieldMetadata.Type
		fieldTag := fieldMetadata.Tag.Get("di") 

		// если нет тэга di, то пропускаем поле
		if fieldTag == "" {
			continue
		} else if fieldTag != ignoreTag && 
			fieldTag != enrichTag && 
			fieldTag != instantiateTag {
			panic(fmt.Sprintf("Невалидный di тэг: %s\n", fieldTag))
		}

		// инициализировать можно только либо указатель, либо интерфейс
		if fieldType.Kind() != reflect.Ptr && fieldType.Kind() != reflect.Interface {
			continue
		}

		// если поле преинициализированно, то пропускаем его
		if !fieldValue.IsNil() {
			continue
		}

		// если поле приватное, то не инициализируем его
		if !fieldValue.CanSet() {
			continue
		}

		if fieldType.Kind() == reflect.Interface {
			// ищем в контейнере имплементацию интерфейса
			// подставляем тип имплементации
			fieldType = c.Resolve(fieldType)
		}

		switch fieldTag {
		// если поле имеет тэг `di:"ignore"`
		// то это поле будет пропущено
		case ignoreTag:
			continue
		// если поле имеет тэг `di:"instantiate"`
		// то будет создан экземпляр структуры,
		// но его поля не будут инициализированы
		case instantiateTag:
			newInstance := reflect.New(fieldType.Elem())
			fieldValue.Set(newInstance)
			if method := newInstance.MethodByName("AfterInstantiated"); 
				method.IsValid() {
				method.Call(nil)
			}

			continue
		// если поле имеет тэг `di:"enrich"`
		// то будет создан экземпляр структуры,
		// и ее поля будут инициализированы
		case enrichTag:
			if singletone, exists := c.GetSingleton(fieldType); exists {
				fieldValue.Set(singletone)
				continue
			} 
			
			newInstance := reflect.New(fieldType.Elem())
			fieldValue.Set(newInstance)

			c.Enrich(newInstance.Interface())
			continue
		}
	}

	if method := reflect.ValueOf(object).MethodByName("AfterEnriched"); 
		method.IsValid() {
		method.Call(nil)
	}
}
