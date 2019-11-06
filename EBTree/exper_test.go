// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package EBTree

import (
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {
	//fmt.Println(getFileSize("/home/mimota/blockbench-master.zip"))

	//fmt.Println(readDir("/home/mimota/ethenv"))
	fmt.Println(ReadDir("/home/mimota/data"))
	//fmt.Println(readDir("/home/mimota/project")

}

func TestWrite(t *testing.T) {
	//fmt.Println(getFileSize("/home/mimota/blockbench-master.zip"))

	//fmt.Println(readDir("/home/mimota/ethenv"))
	content := "1787:29988"
	WriteFile("/home/julia/test.txt", []byte(content))
	//fmt.Println(readDir("/home/mimota/project")
	content = "aaa\nbbbu\n"
	WriteFile("/home/julia/test.txt", []byte(content))

}
func TestGen(t *testing.T) {
	GenTestData(1000)
}

func TestAppend(t *testing.T) {
	AppendToFile("/home/mimota/test.txt", "aaa\n")
}
