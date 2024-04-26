// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//   - if there are no customers, the barber falls asleep in the chair
//   - a customer must wake the barber if he is asleep
//   - if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//     sits in an empty chair if it's available
//   - when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//     and falls asleep if there are none
//   - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//     empty
//   - after the shop is closed and there are no clients left in the waiting area, the barber
//     goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.
package main

import (
	"channel/barber"
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables

var seat_capacity = 4
var arrival_rate = 500                     // milliseconds
var cut_duration = 1000 * time.Millisecond // milliseconds
var time_open = 10 * time.Second

func main() {
	// seed our random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// print welcome message

	color.Yellow("Welcome to Sleeping Barber Problem program")
	color.Yellow("----------------------------------------------")

	// create channels if we need any
	var clientsChan = make(chan string, seat_capacity)
	var doneChan = make(chan bool)

	// create the barbershop
	shop := barber.BarberShop{
		ShopCapacity:     seat_capacity,
		HaircutDuration:  cut_duration,
		NumbersOfBarbers: 0,
		DoneChan:         doneChan,
		ClientsChan:      clientsChan,
		Open:             true,
	}

	color.Green("The shop is open for the day")

	// add barbers
	shop.AddBarber("Messi")
	shop.AddBarber("Temi")
	shop.AddBarber("Latin")
	shop.AddBarber("Pamilerin")
	shop.AddBarber("Omolewa")

	// start the barbershop as a goroutine
	shopClosing := make(chan bool)
	shopClosed := make(chan bool)

	go func() {
		<-time.After(time_open)
		shopClosing <- true
		shop.CloseShopForDay()
		shopClosed <- true
	}()

	// add clients
	i := 1

	go func() {
		for {
			randomArrival := rand.Int() % (2 * arrival_rate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomArrival)):
				shop.AddClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// block until the barbershop is closed
	<-shopClosed

}
