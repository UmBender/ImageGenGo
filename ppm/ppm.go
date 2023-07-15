package ppm

import (
	"fmt"
	"math" 
	"os"
)

type ppm struct {
	name 	string
	height 	uint32
	length 	uint32
	pontos 	[]uint32
}

func Test() {

	my_ppm, _ := Create_ppm("teste_ppm",1080,1920)
	//my_ppm.FillBackGround(0xFF202020)
	//my_ppm.DrawRect(100,200,500,600, 0xFFEFEFEF)

	my_ppm.DrawSoftRect(0,0,1080,1920, 0xFFFF20AA,0xFF202020) //Pretty purple
	//my_ppm.DrawSoftRect(0,0,1080,1920, 0xFFFFAA20,0xFF202020) //Pretty Cyan

	//my_ppm.DrawSoftRect(0,0,1080,1920, 0xFF60FF20,0xFF202020) // Pretty Gree
	//my_ppm.DrawSphere(540,960, 200, 0xFFFF00FF)

	my_ppm.DrawSoftSphere(400, 200,400,0xFFd97aff, 0xFFd920AA)
	my_ppm.DrawSoftSphere(120, 100,732,0xFFd97aff, 0xFFd920AA)
	my_ppm.DrawSoftSphere(200, 600,250,0xFFd97aff, 0xFFd920AA)
	my_ppm.DrawSoftSphere(80, 900,1700,0xFFd97aff, 0xFFd920AA)
	my_ppm.DrawSoftSphere(120, 600,1000,0xFFd97aff, 0xFFd920AA)
	//my_ppm.DrawSphere(400, 600,400,0xFFd97aff)
	my_ppm.DrawSoftSphere(200, 600,250,0xFFd97aff, 0xFFd920AA)
	//my_ppm.DrawSine(0xFFFFFFFF, 0, 1080, 100)
	//my_ppm.DrawSine(0xFFFFFFFF, 0, 300,100)
	
	//my_ppm.DrawSoftSphere(400, 200,400,0xFFd97aff, 0xFFd920AA)
	 

	//my_ppm.DrawLine(10,35,1000,100, 0xFFFFFFFF)
	my_ppm.WritePPM()
}

func (a *ppm) WritePPM() {
	archive, err := os.OpenFile(a.name + ".ppm", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when tring to open image %s", err)
		archive.Close()
		return
	}

	_, err1 :=fmt.Fprintf(archive, "P6\n%d %d\n255\n", a.length, a.height)

	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Error when tring to write image %s", err)
		archive.Close()
		return
	}
	for i := 0; i < int(a.height * a.length); i++ {
		pixel := a.pontos[i]
		pixel &= 0x00FFFFFF
		pixeis := []byte{0,0,0}
		for j := 0; j < 3; j++ {
			pixeis[j] = byte(pixel)
			pixel = (pixel / 256)
		}
		archive.Write(pixeis)
	}
	archive.Close()
}

func Create_ppm(new_name string, new_height uint32, new_length uint32) (ppm, error){
	new_ppm := ppm {
		name: new_name,
		height: new_height,
		length: new_length,
	}
	new_ppm.pontos = make([]uint32, new_height * new_length)
	if new_ppm.pontos != nil {
		return new_ppm, os.ErrDeadlineExceeded
	}
	return new_ppm, nil
}

func (p *ppm) FillBackGround(color uint32) {
	for i := 0; i < int(p.height); i++ {
		for j := 0; j < int(p.length); j++ {
			p.pontos[i * int(p.length) + j] = color
		}
	}
}


func (p *ppm) DrawSphere(i0 int, j0 int, r int, color uint32) {
	for i := i0-r; i < i0+r; i++ {
		for j := j0-r; j < j0+r; j++ {
			jpow := (j-j0)*(j-j0)
			ipow := (i-i0)*(i-i0)
			if ipow + jpow <= r * r {
				p.pontos[i*int(p.length) + j] = color
			}
		}
	}
}



func (p *ppm) DrawRect(i0 int, j0 int, ie int, je int, color uint32)  {
	for i := i0; i < ie; i++{
		for j := j0; j < je; j++ {
			p.pontos[i* int(p.length) + j] = color
		}
	}
}

func Decomp(color uint32) (uint8, uint8, uint8) {
	
	red := uint8(color & 0x0000FF)
	green := uint8((color & 0x00FF00) >> 8)
	blue := uint8((color & 0xFF0000) >> 16)

	return red, green, blue
}

func (p *ppm) DrawSoftSphere(raio uint32, dx uint32, dy uint32, color uint32, backgroundColor uint32) {

	
	red_init, green_init, blue_init := Decomp(backgroundColor)
	grad_cof_red, grad_cof_green, grad_cof_blue :=gradCof(backgroundColor, color, raio)


	for i := 0; i < 2 * int(raio); i++ {
		for j := 0; j < 2 * int(raio); j++ {
			powx := (raio - uint32(i)) * (raio - uint32(i))
			powy := (raio - uint32(j)) * (raio - uint32(j))
			if powx + powy < raio * raio {
				grad_pos_red := (uint32(grad_cof_red * float64(raio* raio - powx - powy)) + uint32(red_init) )
				grad_pos_green := (uint32(grad_cof_green * float64(raio* raio - powx - powy))+ uint32(green_init) )<< 8 
				grad_pos_blue := (uint32(grad_cof_blue * float64(raio* raio - powx - powy))+uint32(blue_init) )<< 16 
				
				
				var grad_pos uint32 = 0xFF000000 | grad_pos_red | grad_pos_green | grad_pos_blue
				
				p.pontos[(i + int(dx)) * int(p.length) +int( dy) + j ] = grad_pos
			}
		}
	}
}


func (p *ppm) DrawSoftRect(i0 int, j0 int, ie int, je int, color uint32,backgroundColor uint32 )  {
	red_init, green_init, blue_init := Decomp(backgroundColor)
	grad_cof_red, grad_cof_green, grad_cof_blue := gradCof(backgroundColor, color, 
	uint32(math.Sqrt((float64(i0)-float64(ie))* (float64(j0)-float64(je)))))

	for i := 0; i < ie - i0; i ++ {
		for j := 0; j < je - j0; j++ {
			grad_pos_red := (uint32(grad_cof_red * float64(i * j)) + (uint32(red_init)))
			grad_pos_green := (uint32(grad_cof_green * float64(i * j)) + (uint32(green_init))) << 8
			grad_pos_blue := (uint32(grad_cof_blue * float64(i * j)) + (uint32(blue_init))) << 16
			var grad_pos uint32 = 0xFF000000 | grad_pos_red | grad_pos_green | grad_pos_blue

			p.pontos[(i + i0) * int(p.length) + int(j0) + j] = grad_pos
		}
	}
}

func (p *ppm) DrawLine(i0 int, j0 int, ie int, je int, color uint32)  {
	a:= (float64(ie) -float64( i0)) / (float64(je) -float64(j0))
	c := float64(ie) - float64(je)* a

	for j := j0; j < je; j++ {
		new_i := int(a * float64(j) + c)
		next_i := int(a * float64(j + 1) +c)

		if j + 1 < je {
			if int(math.Abs(float64(next_i - new_i))) > 1.0 {
				for k := 1; k + new_i < next_i ; k++ {
					p.pontos[(k + new_i) * int(p.length) + j] = color } }
		}
		p.pontos[new_i * int(p.length) + j] = color
	}
}

func gradCof(color_init uint32, color_end uint32, number_levels uint32)  (float64, float64, float64){


	red, green, blue := Decomp(color_end)
	red_init, green_init, blue_init := Decomp(color_init)

	var grad_cof_red float64 = float64(red - red_init) / float64( number_levels * number_levels )
	var grad_cof_green float64 = float64(green - green_init) / float64( number_levels * number_levels )
	var grad_cof_blue float64 = float64(blue - blue_init) / float64( number_levels * number_levels )

	return grad_cof_red, grad_cof_green, grad_cof_blue
}

func (p* ppm) DrawSine(color uint32, desloc_horz int, amplitude int,desloc_vert int)  {
	amplit :=  float64(amplitude)/2.0
	omeg := 2.0 * math.Pi / float64(p.length)
	
	radi_desloc_horz := omeg * float64(desloc_horz)

	
	for j := 0; j < int(p.length) ; j++ {
		i := int((math.Sin(float64(j) * omeg + radi_desloc_horz)) * (amplit - 1)) +int(p.height)/2 -desloc_vert 
		
		if i > 0 && i < int(p.height) {

			p.pontos[i * int(p.length) + j] = color
		}
	}
}

