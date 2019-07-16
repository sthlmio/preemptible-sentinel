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
	"github.com/sthlmio/preemptible-sentinel/pkg/utils"
	"os"
	"time"
)

type Config struct {
	DurationInMinutes time.Duration
	DeleteDiffMinutes float64
}

func Get() Config {
	// Defaults
	c := Config{
		DurationInMinutes: time.Duration(1),
		DeleteDiffMinutes: float64(5),
	}

	if os.Getenv("DURATION_IN_MINUTES") != "" {
		c.DurationInMinutes = time.Duration(utils.StringToInt(os.Getenv("DURATION_IN_MINUTES")))
	}

	if os.Getenv("DELETE_DIFF_MINUTES") != "" {
		c.DeleteDiffMinutes = float64(utils.StringToInt(os.Getenv("DELETE_DIFF_MINUTES")))
	}

	return c
}
