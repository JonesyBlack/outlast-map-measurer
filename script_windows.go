// from Ben Sarah Golightly
package main

import (
	"strconv"
	"syscall"
	"unsafe"

	. "github.com/gonutz/w32/v2"
	"github.com/tawesoft/golib/v2/dialog"
)

type mouseLocation struct {
	x int
	y int
}

func WndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch msg {
	case WM_CLOSE:
		DestroyWindow(hwnd)
		MessageBox(HWND(uintptr(0)), "After the Window has been closed", "Info", MB_ICONINFORMATION|MB_OK)
	case WM_DESTROY:
		PostQuitMessage(0)
	case WM_INPUT:
		in, ok := GetRawInputData(HRAWINPUT(lParam), RID_INPUT)
		if !ok {
			return LRESULT(DefWindowProc(hwnd, msg, uintptr(wParam), uintptr(lParam)))
		}
		if in.Header.Type == RIM_TYPEMOUSE {
			var mouse *RAWMOUSE = (*RAWMOUSE)(unsafe.Pointer(&in))

			var moveMode string
			if (mouse.Flags & MOUSE_MOVE_ABSOLUTE) != 0 {
				moveMode = "abs"
			} else {
				moveMode = "rel"
			}
			str := "mouse: move" + moveMode + "dx=" + strconv.Itoa(int(mouse.LastX)) + " dy=" + strconv.Itoa(int(mouse.LastY))
			MessageBox(HWND(uintptr(0)), str, "Info", MB_ICONINFORMATION|MB_OK)

			if (mouse.ButtonFlags() & uint16(RI_MOUSE_LEFT_BUTTON_DOWN)) != 0 {
				MessageBox(HWND(uintptr(0)), "left down", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_LEFT_BUTTON_UP) != 0 {
				MessageBox(HWND(uintptr(0)), "Left up", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_RIGHT_BUTTON_DOWN) != 0 {
				MessageBox(HWND(uintptr(0)), "Right down", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_RIGHT_BUTTON_UP) != 0 {
				MessageBox(HWND(uintptr(0)), "Right up", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_MIDDLE_BUTTON_DOWN) != 0 {
				MessageBox(HWND(uintptr(0)), "Middle down", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_MIDDLE_BUTTON_UP) != 0 {
				MessageBox(HWND(uintptr(0)), "Middle up", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_4_DOWN) != 0 {
				MessageBox(HWND(uintptr(0)), "x1 down", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_4_UP) != 0 {
				MessageBox(HWND(uintptr(0)), "x1 up", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_5_DOWN) != 0 {
				MessageBox(HWND(uintptr(0)), "x2 down", "Info", MB_ICONINFORMATION|MB_OK)
			}
			if (mouse.ButtonFlags() & RI_MOUSE_BUTTON_5_UP) != 0 {
				MessageBox(HWND(uintptr(0)), "x2 up", "Info", MB_ICONINFORMATION|MB_OK)
			}
		}

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

	var riList []RAWINPUTDEVICE
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
	MessageBox(HWND(0), "After the WNDCLASSEX has been made", "Info", MB_ICONINFORMATION|MB_OK)
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

	// mouseLocations := []mouseLocation{}
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

	// 		mouseLocations = append(mouseLocations, mouseLocation{x: x, y: y})
	// 	}

	// 	if i.Kind == hook.KeyDown && i.Rawcode == 90 {
	// 		dialog.Message{
	// 			Title:  "Info",
	// 			Format: "key down: %d",
	// 		}.WithArgs(i.Rawcode).Raise()
	// 		break
	// 	}
	// }

	// MouseInfoAlert(mouseLocations)
}

func MouseInfoAlert(mLocs []mouseLocation) {
	dialog.Message{
		Title:  "info",
		Format: "useLocations is %v lol",
	}.WithArgs(mLocs).Raise()
}
