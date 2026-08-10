package media

import "testing"

func Test19201080(t *testing.T) {
	ratio := NewAspectRatio(1920, 1080)
	ratio.MakeCommonRatio()
	if ratio.AsString() != "16:9" {
		t.Errorf("Error not correct aspect ratio. Should be 16:9 but is: %s", ratio.AsString())
	}

	ratio = NewAspectRatio(1080, 1920)
	ratio.MakeCommonRatio()
	if ratio.AsString() != "9:16" {
		t.Errorf("Error not correct aspect ratio. Should be 9:16 but is: %s", ratio.AsString())
	}

	ratio = NewAspectRatio(608, 1080)
	ratio.MakeCommonRatio()
	if ratio.AsString() != "9:16" {
		t.Errorf("Error not correct aspect ratio. Should be 9:16 but is: %s", ratio.AsString())
	}

	ratio = NewAspectRatio(100, 100)
	ratio.MakeCommonRatio()
	if ratio.AsString() == "9:16" || ratio.AsString() == "16:9" {
		t.Errorf("Error not correct aspect ratio. Should not be 9:16 or 16:9 but is: %s", ratio.AsString())
	}

}

func TestAspectRatioFromVideoFile(t *testing.T) {
	ratio, err := GetVideoAspectRatio("boots-video-horizontal.mp4")
	if err != nil {
		t.Errorf("Error while getting the videos aspect ratio. : %s", err)
	}

	if ratio != "16:9" {
		t.Errorf("Error not correct aspect ratio. Should be 16:9 but is: %s", ratio)
	}

	ratio, err = GetVideoAspectRatio("boots-video-vertical.mp4")
	if err != nil {
		t.Errorf("Error while getting the videos aspect ratio. : %s", err)
	}

	if ratio != "9:16" {
		t.Errorf("Error not correct aspect ratio. Should be 9:16 but is: %s", ratio)
	}

}
