// from Ben Sarah Golightly
package main

import (
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
	case WM_DESTROY:
		PostQuitMessage(0)
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
	szClassName, _ := syscall.UTF16PtrFromString("myWindowClass")
	szTitleName, _ := syscall.UTF16PtrFromString("My Window Title") // Shouldn't need this once I make it headless, but I just wanna get a window for now because more guides.

	var wc = WNDCLASSEX{
		Style:      0,
		WndProc:    syscall.NewCallback(WndProc),
		ClsExtra:   0,
		WndExtra:   0,
		ClassName:  szClassName,
		Instance:   GetModuleHandle(""),
		Background: (HBRUSH)(COLOR_WINDOW + 1),
	}
	MessageBox(HWND(uintptr(0)), "After the WNDCLASSEX has been made", "Info", MB_ICONINFORMATION|MB_OK)
	var hwnd HWND
	var Msg MSG
	RegisterClassEx(&wc)

	// Step 2: Creating the Window
	hwnd = CreateWindowEx(
		WS_EX_CLIENTEDGE,
		szClassName,
		szTitleName,
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 240, 120,
		HWND(uintptr(0)), HMENU(uintptr(0)), GetModuleHandle(string(0)), unsafe.Pointer(uintptr(0)))

	if hwnd == HWND(uintptr(0)) {
		MessageBox(HWND(uintptr(0)), "Window Creation Failed!", "Error!",
			MB_ICONEXCLAMATION|MB_OK)
		return
	}

	ShowWindow(hwnd, SW_SHOWNORMAL)
	UpdateWindow(hwnd)

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
		Format: "mouseLocations is %v lol",
	}.WithArgs(mLocs).Raise()
}
