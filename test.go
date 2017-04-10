package main


import (
	"fmt"
	"sync/atomic"
	"time"
	"coop-go"
	"flag"
	"os"
	"runtime/pprof"

)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Printf(err.Error())
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

	c1 := int32(0)
	c2 := int32(0)
	count := int32(0)
	
	tickchan := time.Tick(time.Millisecond * time.Duration(1000))
	go func(){
		for{ 
			_ = <- tickchan
			tmp := atomic.LoadInt32(&count)
			atomic.StoreInt32(&count,0)	
			fmt.Printf("count:%d\n",tmp)

		}
	}()


	var p *coop.CoopScheduler


	p = coop.NewCoopScheduler(func (e interface{}){
		atomic.AddInt32(&c1,1)
		atomic.AddInt32(&count,1)
		c2++
		if c1 != c2 {
			fmt.Printf("not equal,%d,%d\n",c1,c2)
		}

		if c2 >= 5000000 {
			p.Close()
			return
		}

		p.Await(func () {
			time.Sleep(time.Millisecond * time.Duration(100))
		})
		
		p.PostEvent(1)
	})

	for i := 0; i < 10000; i++ {
		p.PostEvent(1)
	}


	p.Start()

	fmt.Printf("scheduler stop,total taskCount:%d\n",c2)


}
