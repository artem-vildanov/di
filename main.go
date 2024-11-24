package main

import (
	"di/repos"
	"di/handlers"
	"fmt"
	"reflect"
)

func main() {
	di := NewDependencyContainer()
	
	Register[repos.RepoInterface, repos.RepoImpl](di)

	handler := Instantiate[handlers.UserHandler](di).Handle



	handler()
}

type DependencyContainer struct {
	registry map[reflect.Type]reflect.Type
}

func NewDependencyContainer() *DependencyContainer {
	return &DependencyContainer{
		registry: make(map[reflect.Type]reflect.Type),
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
	if !implType.Implements(reflect.TypeOf((*Injectable)(nil)).Elem()) {
		panic(fmt.Sprintf("Тип %s не реализует интерфейс Injectable", implType))
	}

	container.registry[iType] = implType
}

func (c *DependencyContainer) Resolve(interfaceType reflect.Type) reflect.Type {
	implType, exists := c.registry[interfaceType]
	if !exists {
		panic(fmt.Sprintf("Для интерфейса %s не зарегистрирована реализация", interfaceType))
	}
	return implType
}



type Injectable interface {
	// сработает когда был создан начальный экземпляр структуры, и поля еще не инициализированы
	// можно инициализировать поля начальными данными
	// а также можно инициализировать те поля которые Injectable, тем самым
	// инъекция зависимостей для преинициализированных Injectable полей
	// не будет осуществлена
	AfterInstantiated()
	// сработает после того как инициализированы поля структуры
	AfterEnriched()
}

func Instantiate[T Injectable](container *DependencyContainer) *T {
	var instance T

	if reflect.TypeOf(instance).Kind() == reflect.Interface {
		panic("Нельзя инстанциировать интерфейс")
	}

	instance.AfterInstantiated()
	container.Enrich(&instance)
	instance.AfterEnriched()

	return &instance
}

// обогащаем инстанциированную структуру 
// инициализируя ее поля
// функция вызывается рекурсивно для полей структуры
func (c *DependencyContainer) Enrich(object any) {
	objectValue := reflect.ValueOf(object).Elem()
	objectType := objectValue.Type()

	for i := 0; i < objectType.NumField(); i++ {
		fieldMetadata := objectType.Field(i)
		fieldValue := objectValue.Field(i)
		fieldType := fieldMetadata.Type

		// если поле преинициализированно, то пропускаем его
		if !fieldValue.IsNil() {
			continue
		}

		// если поле приватное, то не инициализируем его
		if !fieldValue.CanSet() {
			continue
		}

		// инициализировать можно только либо указатель, либо интерфейс
		if fieldType.Kind() != reflect.Ptr && fieldType.Kind() != reflect.Interface {
			continue
		}

		if fieldType.Kind() == reflect.Interface {
			// ищем в контейнере имплементацию интерфейса
			// подставляем тип имплементации
			fieldType = c.Resolve(fieldType)
		}

		// если поле имеет тэг `di:"instantiate"`
		// то будет создан экземпляр структуры,
		// но его поля не будут инициализированы
		if fieldMetadata.Tag.Get("di") == "instantiate" {
			newInstance := reflect.New(fieldType.Elem())
			fieldValue.Set(newInstance)
			continue
		}

		// если поле является указателем на структуру, которая не реализует интерфейс Injectable,
		// и у этого поля нет тэга `di:"enrich"`
		// то поле не будет инициализировано
		if !fieldType.Implements(reflect.TypeOf((*Injectable)(nil)).Elem()) &&
			fieldMetadata.Tag.Get("di") != "enrich" {
			continue
		}

		newInstance := reflect.New(fieldType.Elem())
		fieldValue.Set(newInstance)

		callAfterInstantiated(newInstance)
		c.Enrich(newInstance.Interface())
		callAfterEnriched(newInstance)
	}
}

func callAfterInstantiated(instance reflect.Value) {
	afterInstantiatedMethod := instance.MethodByName("AfterInstantiated")
	if afterInstantiatedMethod.IsValid() {
		afterInstantiatedMethod.Call(nil)
	}
}

func callAfterEnriched(instance reflect.Value) {
	afterProvidedMethod := instance.MethodByName("AfterEnriched")
	if afterProvidedMethod.IsValid() {
		afterProvidedMethod.Call(nil)
	}
}

