package patterns

import (
	"fmt"
	"time"
)

/*
pipeline pattern is used when tasks are heavily dependent to one another. It
is somewhat like the loop with a dependency between iterations.

The interesting to note here is that we'll see that even though tasks need
to follow an order (a property not very easy to achieve in a concurrent
system), there is still a speedup to be gained over the sequential execution

The example program would be a recipe:

1. Prepare the baking tray
2. Pour the cupcake mixture
3. Bake the mixture in oven
4. Add toppings
5. Pack cupcakes in a box for delivery

*/

const (
	ovenTime           = 5
	everyThingElseTime = 2
)

// 1. prepare the baking tray function
func PrepareTray(trayNumber int) string {
	fmt.Println("preparing empty tray", trayNumber)
	time.Sleep(everyThingElseTime * time.Second)
	tid := fmt.Sprintf("tray number %d", trayNumber)
	return tid
}

// 2. Pour the cupcake mixture
func MixCupcake(tray string) string {
	fmt.Println("Pouring cupcake Mixture in", tray)
	time.Sleep(everyThingElseTime * time.Second)
	cup := fmt.Sprintf("cupcake in %s", tray)
	return cup
}

// 3. Bake the mixture in oven
func Bake(mixture string) string {
	fmt.Println("baking", mixture)
	time.Sleep(ovenTime * time.Second)
	bakedMix := fmt.Sprintf("baked %s", mixture)
	return bakedMix
}

// 4. Add toppings
func AddTopping(bakedCupcake string) string {
	fmt.Println("Adding topping to", bakedCupcake)
	time.Sleep(everyThingElseTime * time.Second)
	top := fmt.Sprintf("topping on %s", bakedCupcake)
	return top
}

// 5. Pack cupcakes in a box for delivery
func Box(finishedCupcake string) string {
	fmt.Println("Boxing", finishedCupcake)
	time.Sleep(everyThingElseTime * time.Second)
	boxed := fmt.Sprintf("%s is boxed", finishedCupcake)
	return boxed
}

// sequential version of cake
func CakePrepSequentialMain() {
	for i := 0; i < 10; i++ {
		tray := PrepareTray(i)
		mixedCupcake := MixCupcake(tray)
		bakedCupcake := Bake(mixedCupcake)
		doneCupcake := AddTopping(bakedCupcake)
		boxedCupcake := Box(doneCupcake)
		fmt.Println("accepting", boxedCupcake)
	}
}
