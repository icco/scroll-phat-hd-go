package scrollphathd

import (
	"testing"
)

func TestDisplay_Basic(t *testing.T) {
	// Set one pixel - it shouldn't tile, so that should be the only pixel displayed
	dev, disp := getDisplay()
	disp.SetPixel(0, 0, 1)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{1, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})

	disp.Clear()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})
}

func TestDisplay_Fill(t *testing.T) {
	dev, disp := getDisplay()
	disp.Fill(1, 1, 3, 3, 1)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{0, 1, 1},
		{0, 1, 1},
	})
}

func TestDisplay_Scroll(t *testing.T) {
	// Set a pixel outside the frame - buffer should dynamically resize, and nothing
	// should be displayed
	dev, disp := getDisplay()
	disp.SetPixel(3, 3, 1)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})

	// Pixel should move left two and up one, from the bottom right corner
	disp.ScrollTo(2, 1)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{0, 0, 0},
		{0, 1, 0},
	})
}

func TestDisplay_Tile(t *testing.T) {
	dev, disp := getDisplay()
	disp.SetPixel(0, 0, 1)
	disp.SetPixel(1, 1, 2)
	disp.SetPixel(2, 2, 3)
	disp.SetPixel(3, 1, 4)

	// Scroll several tiles forward
	disp.ScrollTo(7, 0)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 1, 0},
		{4, 0, 2},
		{0, 0, 0},
	})

	// Now try with tiling disabled
	dev, disp = getDisplay(scrollphathd.WithTiling(false))
	disp.SetPixel(0, 0, 1)
	disp.SetPixel(1, 1, 2)
	disp.SetPixel(2, 2, 3)
	disp.SetPixel(3, 1, 4)
	disp.ScrollTo(7, 0)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})
}

func TestDisplay_Flip(t *testing.T) {
	dev, disp := getDisplay()
	disp.SetPixel(0, 0, 1)
	disp.SetPixel(1, 1, 2)
	disp.SetPixel(2, 2, 3)
	disp.SetPixel(3, 1, 4)

	// Flip horizontal
	disp.SetFlip(true, false)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 0},
		{4, 0, 2},
		{0, 3, 0},
	})

	// Flip horizontal AND vertical
	disp.SetFlip(true, true)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 3, 0},
		{4, 0, 2},
		{0, 0, 0},
	})

	// Try it with scrolling
	disp.ScrollTo(7, 0)
	disp.Show()
	dev.checkPixels(t, [][]byte{
		{0, 0, 3},
		{0, 4, 0},
		{1, 0, 0},
	})
}

// testDevice is a fake device that implements the Device interface to validate output.
// It's designed to display a 3x3 area.
type testDevice struct {
	buffer [][]byte
}

func (d *testDevice) Width() int                    { return 3 }
func (d *testDevice) Height() int                   { return 3 }
func (d *testDevice) SetBuffer(buffer [][]byte)     { d.buffer = buffer }
func (d *testDevice) SetBrightness(brightness byte) {}
func (d *testDevice) Clear() error                  { return nil }
func (d *testDevice) Show() error                   { return nil }

func getDisplay(opts ...scrollphathd.DisplayOption) (*testDevice, *scrollphathd.Display) {
	dev := &testDevice{}
	return dev, scrollphathd.NewWithDevice(dev, opts...)
}

func (d *testDevice) checkPixels(t *testing.T, expected [][]byte) {
	for y, row := range d.buffer {
		for x, val := range row {
			if val != expected[y][x] {
				t.Fatalf("value at (%d, %d) was different (%d) than expected (%d)", x, y, val, expected[y][x])
			}
		}
	}
}
