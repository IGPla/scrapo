package processor

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGetLinks(t *testing.T) {
	var link1, link2 string
	link1 = "http://testlink1.com"
	link2 = "https://testlink2.com"
	var rawHTML *bytes.Buffer = bytes.NewBufferString(fmt.Sprintf("<html><head></head><body><div><div><a href='%v'>Test link 1></a></div><a href='%v'>Test link 2</a></div></body></html>", link1, link2))
	var links []string = GetHTMLElements(rawHTML, "a", "href")
	if len(links) != 2 {
		t.Errorf("Could not find links. e(%d), g(%d)",
			2,
			len(links))
	}
	if !((links[0] == link1 && links[1] == link2) || (links[0] == link2 && links[1] == link1)) {
		t.Errorf("Could not match expected links. e(%v, %v), g(%v, %v)",
			link1, link2, links[0], links[1])
	}
}

func TestGetImages(t *testing.T) {
	var image1, image2 string
	image1 = "http://img.png"
	image2 = "https://img.jpg"
	var rawHTML *bytes.Buffer = bytes.NewBufferString(fmt.Sprintf("<html><head></head><body><div><div><img src='%v'/></div><img src='%v'/></div></body></html>", image1, image2))
	var images []string = GetHTMLElements(rawHTML, "img", "src")
	if len(images) != 2 {
		t.Errorf("Could not find images. e(%d), g(%d)",
			2,
			len(images))
	}
	if !((images[0] == image1 && images[1] == image2) || (images[0] == image2 && images[1] == image1)) {
		t.Errorf("Could not match expected images. e(%v, %v), g(%v, %v)",
			image1, image2, images[0], images[1])
	}
}
