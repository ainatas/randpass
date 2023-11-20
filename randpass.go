package main

import (
	. "randpass/src"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var window = &MainWindow{}

type meta_pass struct {
	seed int
}

func main() {
	SetWindow()
}

var source_bytelist = [][]byte{
	[]byte("1234567890"),
	[]byte("abcdefghijklmnopqrstuvwxyz"),
	[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	[]byte("!@#_+"),
}

func SetWindow() {
	// logfile, _ := os.OpenFile("log", syscall.O_WRONLY, 777)
	source_byte := BytesToText(source_bytelist)
	byte_in_textedit := source_byte
	byte_list := source_bytelist
	window.Title = "RANDOM PASSWORD"
	window.MinSize = Size{
		Width:  400,
		Height: 600,
	}
	window.MaxSize = Size{
		Width:  400,
		Height: 600,
	}
	window.Size = Size{
		Width:  400,
		Height: 600,
	}
	window.Layout = VBox{}

	var result, history, bytelist, test_textedit *walk.TextEdit
	var check_autocopy, check_fair *walk.CheckBox
	var result_group *walk.Splitter
	var pass_length, pass_count *walk.ComboBox
	text_changed := false
	var pass_wait_copy string

	window.Children = []Widget{
		HSplitter{
			MaxSize: Size{
				Height: 30,
			},
			Children: []Widget{
				Label{
					Text: "字符列表",
					MaxSize: Size{
						Height: 20,
					},
				},
				PushButton{
					Text: "重置",
					OnClicked: func() {
						bytelist.SetText(source_byte)
					},
					MaxSize: Size{
						Width:  30,
						Height: 20,
					},
				},
			},
		},

		TextEdit{
			AssignTo: &bytelist,
			MaxSize: Size{
				Height: 100,
				Width:  60,
			},
			VScroll: true,
			Text:    byte_in_textedit,
			OnTextChanged: func() {
				text_changed = true
			},
		},
		GroupBox{
			MaxSize: Size{
				Width:  400,
				Height: 30,
			},
			Layout: HBox{},
			Children: []Widget{
				CheckBox{
					Text:     "复制",
					AssignTo: &check_autocopy,
					Checked:  true,
					// Enabled:  false,
				},
				CheckBox{
					Text:     "均衡权重",
					AssignTo: &check_fair,
				},
				TextLabel{
					Text: "数量",
				},
				ComboBox{
					Model:        []string{"1", "2", "3", "4", "5"},
					CurrentIndex: 0,
					AssignTo:     &pass_count,
					OnCurrentIndexChanged: func() {
						if pass_count.CurrentIndex() == 0 {
							check_autocopy.SetEnabled(true)
						} else {
							check_autocopy.SetEnabled(false)
						}
					},
				},
				TextLabel{
					Text: "长度",
				},
				ComboBox{
					Model:        []string{"4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"},
					CurrentIndex: 6,
					AssignTo:     &pass_length,
				},
			},
		},
		HSplitter{
			MaxSize: Size{
				Height: 30,
			},
			Children: []Widget{
				PushButton{
					Text: "generate",
					OnClicked: func() {
						pass_wait_copy = ""
						if text_changed {
							byte_in_textedit = bytelist.Text()
							byte_list = TextToBytes(byte_in_textedit)
							text_changed = false
						}
						set_pass_length, err := strconv.Atoi(pass_length.Text())
						if err != nil {
							set_pass_length = 10
						}
						pass_count, err := strconv.Atoi(pass_count.Text())
						if err != nil {
							pass_count = 1
						}
						pass_list := make([][]byte, pass_count)

						seed := time.Now().Nanosecond()

						var wg sync.WaitGroup
						pass_list_tmp := make([]*[]byte, pass_count)
						wg.Add(pass_count)
						for i := 0; i < pass_count; i++ {

							rand_pass := NewPass()
							rand_pass.SetBytes(byte_list)
							rand_pass.SetSeed(seed + i)
							rand_pass.Fair = check_fair.Checked()
							// pass_chan := make(chan []byte)
							go GeneratePassRoutine(&wg, rand_pass, set_pass_length, pass_list_tmp, i)
						}
						wg.Wait()
						for k, pass_pointer := range pass_list_tmp {
							pass_list[k] = *pass_pointer
						}

						// rand_pass := NewPass()
						// rand_pass.SetBytes(byte_list)
						// rand_pass.SetSeed(0)
						// rand_pass.Fair = check_fair.Checked()

						// password := string(rand_pass.MakePass(set_pass_length))
						// pass_wait_copy = string(pass_list[0])
						// add_history_list := pass_list
						add_history_list := make([][]byte, pass_count)
						time := time.Now().Format("2006-01-02 15:04:05")
						for k, pass_byte := range pass_list {
							add_history_list[k] = append([]byte(time+": "), pass_byte...)
						}
						add_history_str := BytesToText(add_history_list)

						pass_str := BytesToText(pass_list)
						history_str := history.Text()
						history_list := TextToBytes(history_str)
						if len(history_list) > 100-len(pass_list) {
							history_list = history_list[:100-len(pass_list)]
						}
						history_str = BytesToText(history_list)
						result.SetText(pass_str)
						history.SetText(add_history_str + "\r\n" + history_str)
						pass_wait_copy = string(pass_list[0])
						if check_autocopy.Checked() && check_autocopy.Enabled() && pass_wait_copy != "" {
							clipboard.WriteAll(pass_wait_copy)
						}
						//tmp_label[0].SetSize(walk.Size{Width: 200, Height: 30})
					},
					MaxSize: Size{
						Width: 60,
					},
				},
				PushButton{
					Text: "copy",
					OnClicked: func() {
						list := TextToBytes(result.Text())
						if string(list[0]) != "" {
							clipboard.WriteAll(pass_wait_copy)
						}

						//tmp_label[0].SetSize(walk.Size{Width: 200, Height: 30})
					},
					MaxSize: Size{
						Width: 60,
					},
				},
				PushButton{
					Text: "记录",
					OnClicked: func() {
						if history.Visible() {
							history.SetVisible(false)
						} else {
							history.SetVisible(true)
						}

						//tmp_label[0].SetSize(walk.Size{Width: 200, Height: 30})
					},
					MaxSize: Size{
						Width: 60,
					},
				},
			},
		},

		VSplitter{
			AssignTo: &result_group,
			MinSize: Size{
				Height: 250,
			},
			Children: []Widget{
				TextEdit{
					ReadOnly: true,
					AssignTo: &result,
					MinSize: Size{
						Height: 100,
						Width:  100,
					},
					Font: Font{
						PointSize: 16,
					},
				},
				TextEdit{
					Visible:  false,
					ReadOnly: true,
					AssignTo: &history,
					MinSize: Size{
						Height: 80,
					},
					Font: Font{
						PointSize: 12,
					},
					VScroll: true,
				},
			},
		},
		VSplitter{
			Visible: false,
			Children: []Widget{
				TextEdit{
					AssignTo: &test_textedit,
					MaxSize: Size{
						Height: 100,
					},
				},
				PushButton{
					Text: "test button",
					OnClicked: func() {
						test_textedit.SetText(strings.ReplaceAll(bytelist.Text(), "\r\n", "\n"))
					},
				},
			},
		},
	}

	window.Run()

}

func GeneratePassRoutine(wg *sync.WaitGroup, r *RandPass, length int, pass_list []*[]byte, index int) {
	password := r.MakePass(length)
	// fmt.Println(password)
	pass_list[index] = &password
	wg.Done()
}

func TextToBytes(t string) [][]byte {
	t = strings.ReplaceAll(strings.ReplaceAll(t, "\r\n", "\n"), "\r", "\n")
	stringlist := strings.Split(t, "\n")
	bytelist := [][]byte{}
	for _, s := range stringlist {
		if len(s) == 0 {
			continue
		}
		bytelist = append(bytelist, []byte(s))

	}
	return bytelist
}

func BytesToText(b [][]byte) string {
	var text string
	for k, bs := range b {
		text += string(bs)
		if k < len(b)-1 {
			text += "\r\n"
		}
	}
	return text
}
