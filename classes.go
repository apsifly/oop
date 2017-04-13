package main

import "github.com/nsf/termbox-go"
import "errors"

type Figure interface {
	Draw(int, int, int, int)
	Intersect(int,int) Figure
	Remove(int, int) error
}

type point struct {
	x int
	y int
	sym rune
}

func NewPoint (x int, y int, sym rune) *point {
	p:=new(point)
	p.x=x
	p.y=y
	p.sym=sym
	return p
}

func (p *point) Draw(x0 int, y0 int, w int, h int) {
	//fmt.Println(string(p.sym))
	termbox.SetCell(p.x, p.y, p.sym, termbox.ColorDefault, termbox.ColorDefault)

}
func (p *point) Intersect(x int, y int) Figure {
	//fmt.Println(string(p.sym))
	if x == p.x && y == p.y {
		return p
	}
	return nil

}
func (p *point) Remove (x int, y int) error {

	return nil //can be removed
}

type field struct {
	w int
	h int
	objects []Figure
}

func NewField() *field {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	f:=new(field)
	f.w, f.h = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	return f
}

func (f *field) Flush() {
	termbox.Flush()
}
func (f *field) Close() {
	termbox.Close()
}
func (f *field) Add(d Figure) {
	f.objects=append(f.objects, d)
}
func (f *field) Redraw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, o := range f.objects {
		o.Draw(0, 0, f.w, f.h)
	}
	termbox.Flush()
}
func (f *field) Intersect (x int, y int) Figure {
	for i,_ := range f.objects {
		j:= f.objects[i].Intersect(x, y)
		if j != nil {
			return j
		}
	}
	return nil
}

func (f *field) Remove (x int, y int) error {
	for i,_ := range f.objects {
		if f.objects[i].Intersect(x, y) != nil {
			err:= f.objects[i].Remove(x, y)
			if err == nil {
				f.objects=append(f.objects[:i],f.objects[i+1:]...)
				return nil
			} else {
				return err
			}
		}
	}
	return errors.New("no object to delete")
}

type snake struct {
	d Direction
	len int
	body []*point
	f *field
}
type Direction int
const (  // iota is reset to 0
        UP Direction = iota  // c0 == 0
        RIGHT Direction = iota  // c1 == 1
        DOWN Direction = iota  // c2 == 2
        LEFT Direction = iota
)

func NewSnake (f *field, x int, y int) *snake {
	s:=new(snake)
	s.d=UP
	s.len=10
	s.f=f
	s.body=append(s.body,NewPoint(x,y,'#'))
	return s
}
func (s *snake) Grow(dir Direction) string {
	if dir == (s.d+2)%4 {
		dir=s.d
	}
	lastx:=s.body[len(s.body)-1].x
	lasty:=s.body[len(s.body)-1].y
	var newx, newy int
	switch dir {
		case 0:
			newx=lastx
			newy=lasty-1

		case 1:
			newx=lastx+1
			newy=lasty

		case 2:
			newx=lastx
			newy=lasty+1

		case 3:
			newx=lastx-1
			newy=lasty

	}
	if s.f.Intersect(newx,newy) != nil {
			 if  s.f.Remove(newx,newy) !=nil {
			 		return "gameover"
			 	} else {
			 		s.len+=1
			 		return "ate"
			 	}
			}
			s.body=append(s.body, NewPoint(newx,newy, '#'))
	s.d=dir
	if len(s.body) > s.len {
		s.body=s.body[1:]
	}
	return ""
}

func (s *snake) Draw (x0 int, y0 int, w int, h int) {
	for _,i := range s.body {
		i.Draw(x0, y0, w, h)
	}
}

func (s *snake) Intersect (x int, y int) Figure {
	for i,_ := range s.body {
		j:= s.body[i].Intersect(x, y)
		if j != nil{
			return j
		}
	}
	return nil
}

func (s *snake) Remove (x int, y int) error {

	return errors.New("will not remove snake")
}

func (s *snake) parseKeyboard(dir *Direction){
		for {
			//l.Println("j cycle")
			ev := termbox.PollEvent()
					if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowUp && s.d != DOWN {
						*dir = UP
					}
					if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowRight && s.d != LEFT {
						*dir = RIGHT
					}
					if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowDown && s.d != UP {
						*dir = DOWN
					}
					if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowLeft && s.d != RIGHT {
						*dir = LEFT
					}
				//l.Println("select")
				//time.Sleep(10 * time.Millisecond)
		}
}

type wall struct {
	body []*point
}

func NewWall(x0 int,y0 int,x1 int,y1 int) *wall {
	w:=new(wall)
	if x0==x1 {
		if y1<y0 { y1,y0=y0,y1 }
		for i:=y0;i<=y1 ; i++ {
			w.body = append(w.body,NewPoint(x0,i, 'O'))
		}
		return w
	}
	if y0==y1 {
		if x1<x0 {x1,x0=x0,x1}
		for i:=x0; i<=x1 ; i++ {
			w.body=append(w.body,NewPoint(i,y0,'O'))
		}
		return w
	}
	panic("attempt to add diagonal line")
	return nil
}
func (wa *wall) Draw(x0 int, y0 int, w int, h int) {
	for _,i := range wa.body {
		i.Draw(x0, y0, w, h)
	}
}
func (w *wall) Intersect (x int, y int) Figure {
	for i,_ := range w.body {
		j:= w.body[i].Intersect(x, y) 
		if j != nil{
			return j
		}
	}
	return nil
}
func (w *wall) Remove (x int, y int) error {

	return errors.New("will not remove wall")
}