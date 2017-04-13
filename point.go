package main


import "time"
//import "log"
//import "os"
import "math/rand"






func main() {
	//l := log.New(os.Stderr, "",1)

	f:= NewField()
	defer f.Close()
	//f.Add(NewPoint(10,2,'g'))
	s:=NewSnake(f, f.w/2, f.h/2)
	f.Add(s)
	f.Add(NewWall(10,10,20,10))
	f.Add(NewWall(0,0,f.w-1,0))
	f.Add(NewWall(0,f.h-1,f.w-1,f.h-1))
	f.Add(NewWall(0,0,0,f.h-1))
	f.Add(NewWall(f.w-1,0,f.w-1,f.h-1))
	for {
		p:=NewPoint(int(rand.Float64()*float64(f.w)),int(rand.Float64()*float64(f.h)),'*')
		if f.Intersect(p.x, p.y) == nil {
			f.Add(p)
			break
		}
	}

	var dir Direction = UP
	go s.parseKeyboard(&dir)
	END:
		for i := 1; i <= 200; i++ {
			r := s.Grow(dir)
			f.Redraw()
			switch r {
			case "":
			case "ate":
				for {
					p:=NewPoint(int(rand.Float64()*float64(f.w)),int(rand.Float64()*float64(f.h)),'*')
					if f.Intersect(p.x, p.y) == nil {
						f.Add(p)
						break
					}
				}
			case "gameover":
				break END
			}

			time.Sleep(300*time.Millisecond)
		}


}