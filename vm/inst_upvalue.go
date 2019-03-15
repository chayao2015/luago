package vm

import (
	. "luago/api"
)

//TODO:
// R(A) := UpValue[B][RK(C)]
func getTabUp(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a++

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)
}
