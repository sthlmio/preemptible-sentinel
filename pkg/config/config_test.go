/*
MIT License

Copyright (c) 2019 sthlm.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package config

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	os.Setenv("DURATION_IN_MINUTES", "10")
	os.Setenv("DELETE_DIFF_MINUTES", "30")

	config := Get()

	if config.DurationInMinutes != 10 {
		t.Errorf("DurationInMinutes was incorrect, got: %d, want: %d.", config.DurationInMinutes, 10)
	}

	if config.DeleteDiffMinutes != float64(30) {
		t.Errorf("DeleteDiffMinutes was incorrect, got: %f, want: %f.", config.DeleteDiffMinutes, float64(30))
	}
}
