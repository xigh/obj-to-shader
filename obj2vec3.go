package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	spaces = regexp.MustCompile("[ \t]+")
	scale  = flag.Float64("scale", 1.0, "scale object")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("usage: obj2vec [name]\n")
		return
	}

	name := flag.Arg(0)

	buf, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		return
	}

	vertices := make([][]float32, 0)
	indices := make([]int, 0)

	lines := strings.Split(string(buf), "\n")
	for n, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		tokens := spaces.Split(line, -1)
		if len(tokens) == 0 {
			continue
		}
		switch tokens[0] {
		case "v":
			if len(tokens) != 4 {
				fmt.Printf("%s:%n: invalid count of vertex floats\n", name, n+1)
				return
			}
			vf := make([]float32, 3)
			for i := 0; i < len(vf); i++ {
				f, err := strconv.ParseFloat(tokens[i+1], 32)
				if err != nil {
					fmt.Printf("%s:%n: invalid float: %v\n", name, n+1, err)
					return
				}
				vf[i] = float32(f) * float32(*scale)
			}
			vertices = append(vertices, vf)

		case "f":
			vi := make([]int, len(tokens)-1)
			for i := 0; i < len(vi); i++ {
				z := strings.Split(tokens[i+1], "/")
				if len(z) != 3 {
					fmt.Printf("%s:%n: invalid index format %q\n", name, n+1, tokens[i+1])
					return
				}
				x, err := strconv.ParseUint(z[0], 10, 32)
				if err != nil {
					fmt.Printf("%s:%n: invalid float: %v\n", name, n+1, err)
					return
				}
				vi[i] = int(x)
			}
			if len(vi) < 3 || len(vi) > 4 {
				fmt.Printf("%s:%n: invalid count of indices\n", name, n+1)
				return
			}
			indices = append(indices, vi[0]-1, vi[1]-1, vi[2]-1)
			if len(vi) == 4 {
				indices = append(indices, vi[0]-1, vi[2]-1, vi[3]-1)
			}

		default:
			continue
		}
	}

	fmt.Printf("const int vCount = %d;\n", len(vertices))
	fmt.Printf("const vec3 vertices[] = vec3[](\n")
	for z, vertex := range vertices {
		comma := ' '
		if z < len(vertices)-1 {
			comma = ','
		}
		fmt.Printf(" vec3(%.5f, %.5f, %.5f)%c\n", vertex[0], vertex[1], vertex[2], comma)
	}
	fmt.Printf(");\n\n")

	num := len(indices) / 3

	fmt.Printf("const int iCount = %d;\n\n", num)
	fmt.Printf("const int indices[] = int[](\n")
	for i := 0; i < num; i++ {
		comma := ' '
		if i < num-1 {
			comma = ','
		}
		fmt.Printf(" %d%c\n", indices[i*3+0]<<20|indices[i*3+1]<<10|indices[i*3+2], comma)
	}
	fmt.Printf(");\n")
}
