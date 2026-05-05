package ui

import (
	"fmt"
	"os"
	"sync/atomic"
)

// Per-frame counter of state-bearing widget constructions, used by the
// debug-mode lifecycle detector. The detector warns when a stateful widget
// is constructed inside the ui.Run build closure (which would reset its
// internal state every frame).
var (
	statefulCtorThisFrame   atomic.Uint64
	statefulLifecycleFrame  atomic.Uint64
	statefulLifecycleWarned atomic.Bool
)

// markStatefulCtor is called by stateful widget constructors so the
// debug-mode lifecycle detector can count constructions per frame. No-op
// when IMMYGO_DEBUG is unset.
func markStatefulCtor() {
	if !debugEnabled {
		return
	}
	statefulCtorThisFrame.Add(1)
}

// debugCheckLifecycle is called after each frame's build + layout. If the
// build closure constructed any stateful widgets after the first couple
// of frames, those widgets reset their state every frame and are likely
// a bug. We warn once.
func debugCheckLifecycle() {
	if !debugEnabled {
		return
	}
	frame := statefulLifecycleFrame.Add(1)
	n := statefulCtorThisFrame.Swap(0)
	// Frames 1 and 2 are legitimate startup; widgets created during the
	// first build pass are real. Past that, any new construction means
	// the user wired a stateful widget into the build closure.
	if frame <= 2 || n == 0 {
		return
	}
	if statefulLifecycleWarned.Swap(true) {
		return
	}
	fmt.Fprintf(os.Stderr,
		"\n[IMMYGO_DEBUG] WARNING: %d stateful widget(s) constructed during frame %d.\n"+
			"Stateful widgets (Dropdown, ContextMenu, Toggle, Slider, Input/Password,\n"+
			"Checkbox, RadioGroup, DatePicker, Accordion, DataGrid, Tree, Drawer, Dialog,\n"+
			"SnackbarManager) must be created once outside ui.Run's build closure and\n"+
			"captured by closure — otherwise their internal state (cursor, scroll, popup\n"+
			"open flag, animation progress) resets every frame and they appear broken.\n\n",
		n, frame)
}
