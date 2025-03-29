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
	ovenTime           = 2
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

/*
add on pipe ties an input channel to an output channel and maps what comes out
of input channel to what is expected by the output channel. We use the quit
channel pattern to break away from the mapping
*/
func AddOnPipe[Input, Output any](quit <-chan int, InToOutMap func(Input) Output, in <-chan Input) chan Output {
	out := make(chan Output)
	go func() {
		defer close(out)
		for {
			select {
			case <-quit:
				return
			case input := (<-in):
				out <- InToOutMap(input)
			}
		}
	}()
	return out
}

// concurrent version of cake
func CakePrepConcurrentMain() {
	input := make(chan int)
	quit := make(chan int)
	trayOutput := AddOnPipe(quit, PrepareTray, input)
	mixCupcakeOutput := AddOnPipe(quit, MixCupcake, trayOutput)
	bakeCupcakeOutput := AddOnPipe(quit, Bake, mixCupcakeOutput)
	doneCupcakeOutput := AddOnPipe(quit, AddTopping, bakeCupcakeOutput)
	output := AddOnPipe(quit, Box, doneCupcakeOutput)

	numInput := 10
	go func() {
		for i := 0; i < numInput; i++ {
			input <- i
		}
	}()

	//
	for j := 0; j < numInput; j++ {
		fmt.Println("received", <-output)
	}
	quit <- 1 // signal quit
}
