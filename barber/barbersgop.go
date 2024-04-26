package barber

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity     int
	HaircutDuration  time.Duration
	NumbersOfBarbers int
	DoneChan         chan bool
	ClientsChan      chan string
	Open             bool
}

func (shop *BarberShop) AddBarber(barber string) {

	shop.NumbersOfBarbers++

	go func() {
		isSleeping := false

		color.Yellow("%s goes to the waiting room to check for a client", barber)

		for {
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is no client in the waiting room so %s takes a nap", barber)
				isSleeping = true
			}

			client, ok := <-shop.ClientsChan
			if ok {
				if isSleeping {
					color.Yellow("%s wakes %s up", client, barber)
					isSleeping = false
				}
				// cut hair
				shop.cutHair(barber, client)
			} else {
				// Send baber home and close the goroutine
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair", barber, client)
	time.Sleep(shop.HaircutDuration)
	color.Green("%s is done cutting %s's hair", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	shop.DoneChan <- true
}

func (shop *BarberShop) CloseShopForDay() {

	color.Cyan("Closing the shop for the day")

	close(shop.ClientsChan)
	shop.Open = false

	for i := 0; i < shop.NumbersOfBarbers; i++ {
		<-shop.DoneChan
	}
	close(shop.DoneChan)

	color.Green("----------------------------------------------------------")
}

func (shop *BarberShop) AddClient(client string) {

	color.Green("*** Client %s arrives", client)

	if shop.Open {

		select {
		case shop.ClientsChan <- client:
			color.Green("*** Client %s takes a seat in the waiting room", client)
		default:
			color.Red("The waiting room is full, so the client %s leaves.", client)
		}

	} else {
		color.Red("The shop is closed so, %s leaves", client)
	}

}
