package process

import (
	"context"
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {
	ctx := context.Background()
	p, err := Start(ctx, "youtube-dl", "https://www.youtube.com/watch?v=N0LZ20ppkNo")
	if err != nil {
		t.Fatal(err)
	}
	for b := range p.Output() {
		fmt.Printf("%q\n", b)
	}
	if err := p.Wait(); err != nil {
		t.Fatal(err)
	}

}
