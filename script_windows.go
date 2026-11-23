// from Ben Sarah Golightly
package main

import (
	"strconv"
	"syscall"

	. "github.com/gonutz/w32/v2"
	"github.com/tawesoft/golib/v2/dialog"
)

type mouseChange struct {
	x    int32
	y    int32
	mode string
}
type location struct {
	x int32
	y int32
}

func (loc1 location) sub(loc2 location) location {
	// does loc2 - loc1
	var loc3 location
	loc3.x = loc2.x - loc1.x
	loc3.y = loc2.y - loc1.y
	return loc3
}

func IgnoreError(ret *uint16, err error) *uint16 {
	return ret
}
func Message(str string) {
	MessageBox(HWND(uintptr(0)), str, "Info", MB_ICONINFORMATION|MB_OK)
}

func AppendPosition(curPos location, delta mouseChange) {
	if delta.mode == "abs" {
		curPos.x = delta.x
		curPos.y = delta.y
	} else { // move mode is relative
		curPos.x += delta.x
		curPos.y += delta.y
	}
}

// I'm assuming that when I get a string of mouse change messages, i will have them
// in an array from oldest to newest, and once the array is filled, call this function.
// I'm also assuming that the conversion between mouse position to camera rotation is
// rot_x = mouse_x * sensitivity * resolution_x * 2pi such that the sensitivity is
// some number low enough so that say a full 360  is like
// x = 250, it's 250 * 1080 * 2pi * (1/(540,000 * pi))
//
//	sensitivity
func CalculateSens(mChanges []mouseChange) location {
	currentPosition := location{0, 0} // assume starting at 0, 0
	for i := 0; i < len(mChanges); i++ {
		AppendPosition(currentPosition, mChanges[i])
	}

	return currentPosition.sub(location{0, 0})
}

var mouseChanges = []mouseChange{}
var textHwnd HWND

func WndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	// Problems:
	// 1. Cannot interact with message boxes
	// 2. Program stops responding very soon after it launches
	// 4. cannot scroll text even though https://stackoverflow.com/questions/15429541/displaying-string-output-in-a-window-using-c-in-win32-api
	// says I should be able to (might actually be related to 1, mayhap. Maybe even 3 is related to 1.)
	switch msg {
	case WM_CREATE:
		textHwnd = CreateWindow(
			IgnoreError(syscall.UTF16PtrFromString("EDIT")),
			IgnoreError(syscall.UTF16PtrFromString("Text Goes Here")), WS_VISIBLE|WS_CHILD|ES_LEFT, 10, 10, 1000, 5000, hwnd, HMENU(0), GetModuleHandle(""), nil)
	case WM_CLOSE:
		DestroyWindow(hwnd)
	case WM_DESTROY:
		PostQuitMessage(0)
	case WM_INPUT:
		in, ok := GetRawInputData(HRAWINPUT(lParam), RID_INPUT)
		if !ok {
			return LRESULT(DefWindowProc(hwnd, msg, uintptr(wParam), uintptr(lParam)))
		}
		if in.Header.Type == RIM_TYPEMOUSE {
			var mouse RAWMOUSE = in.GetMouse()

			var moveMode string
			if (mouse.Flags & MOUSE_MOVE_ABSOLUTE) != 0 {
				moveMode = "abs"
			} else {
				moveMode = "rel"
			}
			str := "mouse: move" + moveMode + "dx=" + strconv.Itoa(int(mouse.LastX)) + " dy=" + strconv.Itoa(int(mouse.LastY))
			mouseChanges = append(mouseChanges, mouseChange{x: mouse.LastX, y: mouse.LastY, mode: moveMode})
			SetWindowText(textHwnd, str)
			// MessageBox(HWND(uintptr(0)), str, "Info", MB_ICONINFORMATION|MB_OK)

			if (mouse.ButtonFlags() & RI_MOUSE_LEFT_BUTTON_DOWN) != 0 {
				go Message("left down")
			}
			if (mouse.ButtonFlags() & RI_MOUSE_LEFT_BUTTON_UP) != 0 {
				go Message("Left up")
			}
			if (mouse.ButtonFlags() & RI_MOUSE_RIGHT_BUTTON_DOWN) != 0 {
				go Message("Right down")
				positionChange := CalculateSens(mouseChanges)
				str := strconv.Itoa(int(positionChange.x))
				Message(str)
				DestroyWindow(hwnd)
				// TODO: MAKE THE THING ACTUALLY BE WHEN YOU PRESS Z AND NOT RIGHT CLICK, ALSO, FIX THE THING TO NOT USE THE DIALOG IMPORT NOW
			}
			if (mouse.ButtonFlags() & RI_MOUSE_RIGHT_BUTTON_UP) != 0 {
				go Message("Right up")
			}
			// if (mouse.ButtonFlags() & RI_MOUSE_MIDDLE_BUTTON_DOWN) != 0 {
			// 	MessageBox(HWND(uintptr(0)), "Middle down", "Info", MB_ICONINFORMATION|MB_OK)
			// }
			// if (mouse.ButtonFlags() & RI_MOUSE_MIDDLE_BUTTON_UP) != 0 {
			// 	MessageBox(HWND(uintptr(0)), "Middle up", "Info", MB_ICONINFORMATION|MB_OK)
			// }
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_4_DOWN) != 0 {
				go Message("x1 down")
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_4_UP) != 0 {
				go Message("x1 up")
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_5_DOWN) != 0 {
				go Message("x2 down")
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_5_UP) != 0 {
				go Message("x2 up")
			}
		}
		return LRESULT(DefWindowProc(hwnd, msg, uintptr(wParam), uintptr(lParam)))

	default:
		return LRESULT(DefWindowProc(hwnd, msg, uintptr(wParam), uintptr(lParam)))
	}
	return 0
}

//Step 1: Registering the Window Class

//	func OnMouseMoveRaw(x int, y int) {
//		this->eventBuffer.push(MouseEvent(MouseEvent::EventType::RAW_MOVE, x, y))
//	}
func main() {
	szClassName, _ := syscall.UTF16PtrFromString("myTrialsMeasurer")
	szTitleName, _ := syscall.UTF16PtrFromString("TrialsMeasurer") // Shouldn't need this once I make it headless, but I just wanna get a window for now because more guides.

	var riList [2]RAWINPUTDEVICE
	riList[0].UsagePage = 0x05 // HID_USAGE_GAME_CONTROLS
	riList[0].Usage = 20       // HID_USAGE_Point of View
	riList[0].Flags = 0        // adds game pad
	riList[0].Target = 0

	riList[1].UsagePage = 0x01       // HID_USAGE_PAGE_GENERIC
	riList[1].Usage = 0x02           // HID_USAGE_GENERIC_MOUSE
	riList[1].Flags = RIDEV_NOLEGACY // adds mouse and also ignores legacy mouse messages
	riList[1].Target = 0

	if RegisterRawInputDevices(riList[0], riList[1]) {
		//registration failed. Call GetLastError for the cause of the error
	}

	var wc = WNDCLASSEX{
		Style:      0,
		WndProc:    syscall.NewCallback(WndProc),
		ClsExtra:   0,
		WndExtra:   0,
		ClassName:  szClassName,
		Instance:   GetModuleHandle(""),
		Background: (HBRUSH)(COLOR_WINDOW + 1),
	}
	// MessageBox(HWND(0), "After the WNDCLASSEX has been made", "Info", MB_ICONINFORMATION|MB_OK)
	var hwnd HWND
	var Msg MSG
	if RegisterClassEx(&wc) == 0 {
		MessageBox(HWND(0), "Window Registration Failed!", "Error!",
			MB_ICONEXCLAMATION|MB_OK)
		return
	}

	// Step 2: Creating the Window
	hwnd = CreateWindowEx(
		WS_EX_CLIENTEDGE,
		szClassName,
		szTitleName,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 240, 120,
		HWND(uintptr(0)), HMENU(uintptr(0)), GetModuleHandle(""), nil)

	if hwnd == HWND(uintptr(0)) {
		MessageBox(HWND(uintptr(0)), "Window Creation Failed!", "Error!",
			MB_ICONEXCLAMATION|MB_OK)
		return
	}

	ShowWindow(hwnd, SW_SHOWNORMAL)
	UpdateWindow(hwnd)
	MessageBox(HWND(uintptr(0)), "After the window has been updated for the first time, HWND: "+strconv.FormatUint(uint64(uintptr(hwnd)), 10), "Info", MB_ICONINFORMATION|MB_OK)

	for GetMessage(&Msg, HWND(uintptr(0)), 0, 0) > 0 {
		TranslateMessage(&Msg)
		DispatchMessage(&Msg)
	}

	// qYesNo, _ := dialog.Message{
	// 	Title:  "result",
	// 	Format: "That wall is %d gridspaces long. Would you like to copy a snippet of the json code with the exact length?",
	// }.WithArgs(4).Ask()
	// if qYesNo {
	// 	clipboard.WriteAll("Hi how are ya doing on this fine wonderful day of days?")
	// }
	// hk := hook.Start()
	// defer hook.End()

	// for {
	// 	i := <-hk

	// 	if i.Kind == hook.MouseMove {
	// 		x, y := robotgo.Location()

	// 	}

	// 	if i.Kind == hook.KeyDown && i.Rawcode == 90 {
	// 		dialog.Message{
	// 			Title:  "Info",
	// 			Format: "key down: %d",
	// 		}.WithArgs(i.Rawcode).Raise()
	// 		break
	// 	}
	// }
}

func MouseInfoAlert(mLocs []mouseChange) {
	dialog.Message{
		Title:  "info",
		Format: "useLocations is %v lol",
	}.WithArgs(mLocs).Raise()
}
