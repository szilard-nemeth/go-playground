package main

import (
	"fmt"
	"reflect"
)

type WorkflowHandlerComms struct {
	CleanupChannels map[string]map[Component]chan struct{} //key: cluster id, value: map (key: component, value: channel)
	ResultChannels  map[string]map[Component]chan string   //key: cluster id, value: map (key: component, value: channel)
}

type ChartRev int

type Status string

type Ctx struct {
	ClusterId string
	Comms     *WorkflowHandlerComms
}

func NewWorkflowHandlerComms() *WorkflowHandlerComms {
	return &WorkflowHandlerComms{
		CleanupChannels: make(map[string]map[Component]chan struct{}),
		ResultChannels:  make(map[string]map[Component]chan string),
	}
}

func NewCtx(clusterId string) *Ctx {
	return &Ctx{ClusterId: clusterId, Comms: NewWorkflowHandlerComms()}
}

type Component interface {
	Start(ctx *Ctx) error
	Status(ctx *Ctx, varName string) (Status, error)
}

type ArbitraryComponent struct {
	ComponentName       string
	CurrentChartVersion ChartRev
	DeploymentName      string
}

func (c *ArbitraryComponent) Start(ctx *Ctx) error {
	if _, ok := ctx.Comms.CleanupChannels[ctx.ClusterId]; !ok {
		ctx.Comms.CleanupChannels[ctx.ClusterId] = make(map[Component]chan struct{})
	}
	if _, ok := ctx.Comms.ResultChannels[ctx.ClusterId]; !ok {
		ctx.Comms.ResultChannels[ctx.ClusterId] = make(map[Component]chan string)
	}

	if _, ok := ctx.Comms.CleanupChannels[ctx.ClusterId][c]; !ok {
		ctx.Comms.CleanupChannels[ctx.ClusterId][c] = make(chan struct{})
		ctx.Comms.ResultChannels[ctx.ClusterId][c] = make(chan string, 3)
	}

	//Let's push some data to the channel
	select {
	case ctx.Comms.ResultChannels[ctx.ClusterId][c] <- "ready":
		fmt.Println("ready pushed to channel")
	default:
		fmt.Println("Data NOT pushed to channel")
	}
	ch := ctx.Comms.ResultChannels[ctx.ClusterId][c]
	fmt.Printf("component: %s, result channel data: ch: %s (%v), len(ch): %s\n", c.ComponentName, ch, ch, len(ch))

	return nil
}

func (c *ArbitraryComponent) Status(ctx *Ctx, varName string) (Status, error) {
	ch := ctx.Comms.ResultChannels[ctx.ClusterId][c]
	fmt.Printf("Status. component: %s (varname: %s), Channel data: ch: %s (%v), len(ch: %d)\n", c.ComponentName, varName, ch, ch, len(ch))

	keys := make([]any, len(ctx.Comms.ResultChannels[ctx.ClusterId]))
	i := 0
	for k := range ctx.Comms.ResultChannels[ctx.ClusterId] {
		keys[i] = k
		i++
	}
	fmt.Printf("Status. component: %s (varname: %s), reflect.TypeOf(keys[0]).String(): %s, reflect.TypeOf(keys[0]).Kind(): %s\n",
		c.ComponentName, varName, reflect.TypeOf(keys[0]).String(), reflect.TypeOf(keys[0]).Kind())
	fmt.Printf("Status. component: %s (varname: %s), reflect.TypeOf(c).String(): %s, reflect.TypeOf(c).Kind(): %s\n",
		c.ComponentName, varName, reflect.TypeOf(c).String(), reflect.TypeOf(c).Kind())
	fmt.Printf("Status. component: %s (varname: %s), reflect.TypeOf(*c).String(): %s, reflect.TypeOf(*c).Kind(): %s\n",
		c.ComponentName, varName, reflect.TypeOf(*c).String(), reflect.TypeOf(*c).Kind())

	return "OK", nil
}

func main() {
	clusterId := "cluster111"
	ctx := NewCtx(clusterId)
	c1_1 := ArbitraryComponent{
		ComponentName:       "c1",
		CurrentChartVersion: 1,
		DeploymentName:      "c1d",
	}

	//Later on component fetched from DB and created again
	//object is new, but fields are the same so the hash of the object in the Map should be the same
	c1_2 := ArbitraryComponent{
		ComponentName:       "c1",
		CurrentChartVersion: 1,
		DeploymentName:      "c1d",
	}
	c1_1.Start(ctx)
	c1_1.Status(ctx, "c1_1")
	fmt.Println()
	fmt.Println()
	c1_2.Status(ctx, "c1_2")

	//OUTPUT:
	//ready pushed to channel
	//component: c1, result channel data: ch: %!s(chan string=0x14000124120) (0x14000124120), len(ch): %!s(int=1)
	//Status. component: c1 (varname: c1_1), Channel data: ch: %!s(chan string=0x14000124120) (0x14000124120), len(ch: 1)
	//Status. component: c1 (varname: c1_1), reflect.TypeOf(keys[0]).String(): *main.ArbitraryComponent, reflect.TypeOf(keys[0]).Kind(): ptr
	//Status. component: c1 (varname: c1_1), reflect.TypeOf(c).String(): *main.ArbitraryComponent, reflect.TypeOf(c).Kind(): ptr
	//Status. component: c1 (varname: c1_1), reflect.TypeOf(*c).String(): main.ArbitraryComponent, reflect.TypeOf(*c).Kind(): struct
	//
	//
	//Status. component: c1 (varname: c1_2), Channel data: ch: %!s(chan string=<nil>) (<nil>), len(ch: 0)
	//Status. component: c1 (varname: c1_2), reflect.TypeOf(keys[0]).String(): *main.ArbitraryComponent, reflect.TypeOf(keys[0]).Kind(): ptr
	//Status. component: c1 (varname: c1_2), reflect.TypeOf(c).String(): *main.ArbitraryComponent, reflect.TypeOf(c).Kind(): ptr
	//Status. component: c1 (varname: c1_2), reflect.TypeOf(*c).String(): main.ArbitraryComponent, reflect.TypeOf(*c).Kind(): struct

}
