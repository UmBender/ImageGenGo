package ppm

import (
	"fmt"
	"os"
)

type ppm struct {
	name 	string
	height 	uint32
	length 	uint32
	pontos 	[]uint32
}
const height = 1080
const length = 1920
const background = 0x00FF2020
const red = 0xFF2020FF
const green = 0xFF20FF20
const blue = 0xFFFF2020
const purple = 0xFFFF20FF


var pontos [height * length]uint32

func Test() {
	my_ppm, _ := Create_ppm("new_ppm",1080,1920)
	my_ppm.FillBackGround(0xFF202020)
	my_ppm.DrawRect(100,200,500,600, 0xFFEFEFEF)
	my_ppm.DrawSphere(540,960, 100, 0xFFFF00FF)
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
			p.pontos[i*length + j] = color
		}
	}
}

func Decomp(color uint32) (uint8, uint8, uint8) {
	
	red := uint8(color & 0x0000FF)
	green := uint8((color & 0x00FF00) >> 8)
	blue := uint8((color & 0xFF0000) >> 16)

	return red, green, blue
}

func PintaFundo() {
	for i:= 0; i < height; i++{
		for j := 0; j < length; j++ {
			pontos[i * length + j] =  background
		}
	}
}

func PintaBolinha() {
	raio := 300
	dx := 540 - raio
	dy := 960 - raio
	red_purple, green_purple, blue_purple:= Decomp(purple)

	red_init, green_init, blue_init := Decomp(background)

	var grad_cof_red float64 = float64(red_purple - red_init) / float64(raio * raio)
	var grad_cof_green float64 = float64(green_purple - green_init) / float64(raio * raio)
	var grad_cof_blue float64 = float64(blue_purple - blue_init) / float64(raio * raio)
	for i := 0 ; i < 2*raio; i++ {
		for j:= 0; j < 2 * raio; j++ {
			powx := (raio - i) * (raio - i)
			powy := (raio - j) * (raio - j)
			if powx + powy <= raio * raio {
				grad_pos_red := float64(raio * raio) -float64(powx)  -float64(powy) 
				grad_pos_green := float64(raio * raio) -float64(powx)  -float64(powy) 
				grad_pos_blue := float64(raio * raio) -float64(powx)  -float64(powy) 
				grad_pos := ((uint32(grad_pos_red * grad_cof_red)+ uint32(red_init)) << 0)|
				((uint32(grad_pos_green * grad_cof_green)+ uint32(green_init) ) << 8)|
				((uint32(grad_pos_blue * grad_cof_blue)+  uint32(blue_init) ) << 16)
				pontos[(i + dx) * length + dy + j ] = grad_pos
			}

		}

	}
}
func EscreveArquivo(nome string) {
	arquivo, err := os.OpenFile(nome + ".ppm", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Erro ao abrir:", err)
	}

	
	_, err1 := fmt.Fprintf(arquivo,"P6\n%d %d\n255\n",length, height)
	if err1 != nil {
		return
	}
	for i:= 0; i < height * length; i++ {
		pixel := pontos[i]
		pixel &= 0x00FFFFFF
		pixeis := []byte{0,0,0}
		for j:= 0; j < 3; j++ {
			pixeis[j] = byte(pixel)
			pixel /= 256
		}
		arquivo.Write(pixeis)
	}
	arquivo.Close()
}
