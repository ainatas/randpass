package src

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

var default_bytes = [][]byte{
	[]byte("1234567890"),
	[]byte("abcdefghijklmnopqrstuvwxyz"),
	[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	[]byte("!@#_+"),
}

type RandPass struct {
	Bytelist     [][]byte
	Fair         bool
	seed         int
	bytelist_set bool
}

func NewPass() *RandPass {
	return &RandPass{
		Bytelist:     default_bytes,
		seed:         time.Now().Nanosecond(),
		bytelist_set: true,
	}

}

func (p *RandPass) SetSeed(seed int) {
	if seed == 0 {
		p.seed = time.Now().Nanosecond()
	} else {
		p.seed = seed
	}

}

func (p *RandPass) GetBytesFromFile(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		p.Bytelist = default_bytes
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) > 0 {
			p.Bytelist = append(p.Bytelist, scanner.Bytes())
		}

	}
	p.bytelist_set = true
	return nil
}
func (p *RandPass) SetBytes(bytes [][]byte) {
	p.Bytelist = bytes
	p.bytelist_set = true
}

func (p *RandPass) SetFair() {
	p.Fair = true
}

func (p *RandPass) MakePass(length int) []byte {
	if !p.bytelist_set {
		p.Bytelist = default_bytes
	}
	bytes := p.Bytelist
	line_num := len(bytes)
	if length < line_num {
		length = line_num
	}
	if p.seed == 0 {
		p.SetSeed(0)
	}
	seed := p.seed
	rand := rand.New(rand.NewSource(int64(seed)))
	tmp_byte_list := make([]byte, length)
	all_byte_list := []byte{}
	for k, line := range bytes {
		len_line := len(line)
		index := rand.Intn(len_line)
		tmp_byte_list[k] = line[index]
		all_byte_list = append(all_byte_list, line...)
	}
	if p.Fair {
		for i := line_num; i < length; i++ {
			line_index := rand.Intn(line_num)
			tmp_byte_num := len(bytes[line_index])
			byte_index := rand.Intn(tmp_byte_num)
			tmp_byte_list[i] = bytes[line_index][byte_index]
		}
	} else {
		all_byte_len := len(all_byte_list)
		for i := line_num; i < length; i++ {
			byte_index := rand.Intn(all_byte_len)
			tmp_byte_list[i] = all_byte_list[byte_index]
		}
	}
	byte_list := make([]byte, length)
	tmp_rand_list := NoRepeatRand(length, seed)
	for k, index := range tmp_rand_list {
		byte_list[k] = tmp_byte_list[index]
	}
	return byte_list
}

// generate random number list without repeat begin 0
func NoRepeatRand(length int, seed int) []int {
	tmp_list := make([]int, length)
	rand_list := make([]int, length)
	for i := 0; i < length; i++ {
		tmp_list[i] = i
	}
	rand := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < length-1; i++ {
		tmp_len := len(tmp_list)
		tmp_index := rand.Intn(tmp_len)
		rand_list[i] = tmp_list[tmp_index]
		tmp_list = append(tmp_list[:tmp_index], tmp_list[tmp_index+1:]...)
	}
	rand_list[length-1] = tmp_list[0]
	return rand_list
}
