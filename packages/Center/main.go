package Center

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/iskaa02/qalam/gradient"
	"golang.org/x/term"
)

type Terminal struct {
	Body       string
	BodyMiddle string
	List       []string
	list_num   int
	RGB        []string
	Grad       gradient.Gradient
	Head       string
}

func init() {
	//Clear()
	fmt.Print("\033[s")
}

func App(rgb ...string) Terminal {
	T := Terminal{RGB: rgb}
	if len(rgb) != 0 {
		T.Grad = T.InitRGB()
	}
	return T
}

func (App *Terminal) InitRGB() gradient.Gradient {
	g, _ := gradient.NewGradientBuilder().
		HtmlColors(App.RGB...).
		Mode(gradient.BlendRgb).
		Build()
	return g
}

func (Data *Terminal) Header(s string) {
	if len(Data.RGB) != 0 {
		s = Data.Grad.Mutline(s)
	}
	Data.Head = s
	Data.Body += Data.Head
}
func (Data *Terminal) Print(content string) {
	info := Center(content)
	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	Data.Body += info
	fmt.Print(Data.Body)
}

func (Data *Terminal) PrintUncached(content string, clear bool) {
	info := Center(content)
	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	if clear {
		Clear()
	}
	fmt.Print(Data.Body + "\n" + info)
}

func (Data *Terminal) PrintMiddle(content string) {
	info := Center(content)

	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	Data.BodyMiddle += info
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	Clear()
	fmt.Print(strings.Repeat("\n", (h)/3), Data.BodyMiddle, strings.Repeat("\n", (h/3)))
}

func (Data *Terminal) PrintMiddleUncachedToBody(content string) {
	info := Center(content)
	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	Clear()
	fmt.Print(strings.Repeat("\n", (h)/3), info, strings.Repeat("\n", (h/3)))
}

func (Data *Terminal) InputMiddle(Prefix string) string {
	Copy := Prefix
	Prefix = strings.Replace(Prefix, " ", "", -1)
	info := Center(Prefix)

	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	Clear()
	fmt.Print(strings.Repeat("\n", (h)/2), Data.BodyMiddle, info[:strings.Index(info, Prefix[len(Prefix)-1:])+1]+Copy[strings.Index(Copy, Prefix[len(Prefix)-1:])+1:])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (Data *Terminal) AddList(item string) {
	Data.List = append(Data.List, item)
	Clear()
	fmt.Print(Data.Body)

	if len(Data.List) >= 11 {
		for _, item := range Data.List[Data.list_num-9 : Data.list_num] {
			info := Center(item)
			if len(Data.RGB) != 0 {
				info = Data.Grad.Mutline(info)
			}
			fmt.Print(info)
		}
	} else {
		if Data.list_num > 0 {
			for _, item := range Data.List[0:Data.list_num] {
				info := Center(item)
				if len(Data.RGB) != 0 {
					info = Data.Grad.Mutline(info)
				}
				fmt.Print(info)
			}
		}
	}

	Data.list_num++

}

func (Data *Terminal) ClearList() {
	Clear()
	fmt.Print(Data.Body)
}

func (Data *Terminal) Input(Prefix string) string {
	Copy := Prefix
	Prefix = strings.Replace(Prefix, " ", "", -1)
	info := Center(Prefix)
	if len(Data.RGB) != 0 {
		info = Data.Grad.Mutline(info)
	}
	Clear()
	fmt.Print(info[:strings.Index(info, Prefix[len(Prefix)-1:])+1] + Copy[strings.Index(Copy, Prefix[len(Prefix)-1:])+1:])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func Center(s string) (values string) {
	// get screen size
	w, _, _ := term.GetSize(int(os.Stdout.Fd()))

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		w_ := scanner.Text()
		values += strings.Repeat(" ", w/2-(len(w_)/2)) + w_ + "\r\n"
	}
	return
}

func BottomLeft(s string) {
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Print(strings.Repeat("\n", (h-2)) + s)
}

func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
