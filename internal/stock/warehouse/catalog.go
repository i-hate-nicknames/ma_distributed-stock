package warehouse

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

// todo: maybe create a global type for map[string][]int64?
// this is getting on me nerves

// Catalog maps warehouse address to its items
// Each warehouse data is valid when it's added to the catalog
// but may later get invalid, e.g. when a warehouse goes down or
// somehow loses items
type Catalog struct {
	mux        sync.Mutex
	warehouses map[string][]int64
}

// MakeCatalog makes a single instance of an address book
func MakeCatalog() *Catalog {
	warehouses := make(map[string][]int64, 0)
	return &Catalog{warehouses: warehouses}
}

// GetWarehouses returns all warehouses in this catalog
func (c *Catalog) GetWarehouses() map[string][]int64 {
	return c.warehouses
}

// Copy returs a deep copy of this catalog, which may
// be safely modified
func (c *Catalog) Copy() *Catalog {
	res := MakeCatalog()
	for addr, items := range c.warehouses {
		new := make([]int64, len(items))
		copy(new, items)
		res.warehouses[addr] = new
	}
	return res
}

func (c *Catalog) ApplyShipment(shipment map[string][]int64) {

}

// this doesn't perform the request but only removes item
// from the catalog
func (c *Catalog) PopItem(address string) error {
	wh, ok := c.warehouses[address]
	if !ok {
		return fmt.Errorf("warehouse %s not found", address)
	}
	if len(wh) == 0 {
		return fmt.Errorf("warehouse %s is empty", address)
	}
	c.warehouses[address] = wh[1:]
	return nil
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) AddWarehouse(address string, items []int64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.warehouses[address] = items
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) HasWarehouse(address string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	return ok
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) RemoveWarehouse(address string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	if ok {
		delete(c.warehouses, address)
	}
}

// todo: separate catalogue manipulation code from warehouse network requests
// and from warehouse discovery

// GetInvitations starts listening for invitation messages that active
// warehouses send over UDP. Whenever it gets new warehouse address it sends it
// to the addresses channel
func GetInvitations(ctx context.Context, addresses chan<- string) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			log.Println(err)
		}
		addresses <- string(buf[:n])
		time.Sleep(discoverTimeout)
	}
}

// TakeItems simulates taking items: send take item requests to all available warehouses
func TakeItems(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	toTake := []int64{1, 2}
	for addr := range catalog.warehouses {
		log.Printf("Taking %v from %s\n", toTake, addr)
		ctx := context.Background()
		err := doTakeItems(ctx, addr, toTake)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

// GreetWarehouses sends greeting to every warehouse to test connection
func GreetWarehouses(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	for addr := range catalog.warehouses {
		ctx := context.Background()
		doHello(ctx, addr)
	}
}
