// MIT License
//
// Copyright (c) 2022 Spiral Scout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package container

import endure "github.com/roadrunner-server/endure/pkg/container"

// NewContainer creates endure container with all required options (based on container Config). Logger is nil by
// default.
func NewContainer(cfg Config) (*endure.Endure, error) {
	endureOptions := []endure.Options{
		endure.GracefulShutdownTimeout(cfg.GracePeriod),
	}

	if cfg.PrintGraph {
		endureOptions = append(endureOptions, endure.Visualize(endure.StdOut, ""))
	}

	return endure.NewContainer(cfg.Logger, endureOptions...)
}