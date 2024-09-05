package main

import (
	"fmt"
	sysLog "log"
	"math"
	"unsafe"

	"github.com/rs/zerolog/log"
)

// Define the Shape interface
type Shape interface {
	Area() float64
}

// Rectangle struct
type Rectangle struct {
	Width  float64
	Height float64
}

// Circle struct
type Circle struct {
	Radius float64
}

// Implement the Area method for Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Implement the Area method for Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Function to print the area of any Shape
func PrintArea(s Shape) {
	// fmt.Printf("Area: %.2f\n", s.Area())
	addr := (*Shape)(Noescape2(unsafe.Pointer(&s)))
	log.Printf("addr of a in bar = %v\n", (*int)(Noescape2(unsafe.Pointer(&s))))
	log.Printf("addr of a in bar3333 = %v\n", *addr)
}

func arrayAny(a []any) {
	aaaaaa := a
	_ = (*[]any)(Noescape2(unsafe.Pointer(&aaaaaa)))
	for _, p := range a {
		fmt.Printf("%+v\n", p)
	}
	log.Printf("arrayAny = %v")
	// log.Printf("arrayAny = %v", aaaaaa)
}

func main1() {
	rectangle := Rectangle{Width: 5, Height: 3}
	rectangle2 := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 2.5}
	// Call PrintArea on rectangle and circle, both of which implement the Shape interface
	PrintArea(rectangle) // Prints the area of the rectangle
	// PrintArea(circle)    // Prints the area of the circle

	log.Printf("addr of a in bar5555 = %v\n", circle)

	// arrayAny([]any{&circle, &rectangle})

	sysLog.Printf("", rectangle2)
}
