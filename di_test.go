package di

import "testing"

type testingStruct1 struct {
	struct2 *testingStruct2
	struct3 *testingStruct3
}

func (t *testingStruct1) Construct(struct2 *testingStruct2, struct3 *testingStruct3) {
	t.struct2 = struct2
	t.struct3 = struct3
}

func (t *testingStruct1) testingStruct1Func() string {
	testingStruct2Res := t.struct2.testingStruct2Func()
	testingStruct3Res := t.struct3.testingStruct3Func()
	return "I am testingStruct1Func\n" + testingStruct2Res + testingStruct3Res
} 

type testingStruct2 struct {
	interface4 testingInterface4
}

func (t *testingStruct2) Construct(impl4 testingInterface4) {
	t.interface4 = impl4
}

func (t *testingStruct2) testingStruct2Func() string {
	testingInterface4Res := t.interface4.doSmth()
	return "I am testingStruct2Func\n" + testingInterface4Res
} 

type testingStruct3 struct {}

func (t *testingStruct3) Construct() {
	// empty constructor
}

func (t *testingStruct3) testingStruct3Func() string {
	return "I am testingStruct3Func\n"
} 

type testingInterface4 interface {
	doSmth() string
}

type testingStruct4 struct {
	struct5 *testingStruct5
}

func (t *testingStruct4) Construct(struct5 *testingStruct5) {
	t.struct5 = struct5
}

func (t *testingStruct4) doSmth() string {
	struct5Res := t.struct5.testingStruct5Func() 
	return "I am doSmth\n" + struct5Res
}

type testingStruct5 struct {}

func (t *testingStruct5) testingStruct5Func() string {
	return "I am testingStruct5Func\n"
} 

func TestPositive(t *testing.T) {
	di := NewDependencyContainer()
	Bind[testingInterface4, testingStruct4](di)

	testingStruct1Inst := InitHandler[testingStruct1](di)

	actual := testingStruct1Inst.testingStruct1Func()
	expected := "I am testingStruct1Func\nI am testingStruct2Func\nI am doSmth\nI am testingStruct5Func\nI am testingStruct3Func\n"
		
	if actual != expected {
		t.Errorf("\nExpected:\n%s\nActual:\n%s\n", expected, actual)
	}
}