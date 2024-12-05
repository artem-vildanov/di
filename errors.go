package di

import "fmt"

func panicBindingNotFound(interfaceName string) {
	panic(fmt.Sprintf("Для интерфейса %s не зарегистрирована реализация", interfaceName))
}

func panicFirstArgumentNotInterface() {
	panic("Первый аргумент должен быть интерфейсом")
}

func panicTypeDoesntImplementInterface(typeName string, interfaceName string) {
	panic(fmt.Sprintf("Тип %s не реализует интерфейс %s", typeName, interfaceName))
}

func panicHandlerShouldntBeAnInterface() {
	panic("Handler не может быть интерфейсом")
}

func panicTypeCastFailed() {
	panic("Не удалось произвести приведение типов при инициализации Handler")
}