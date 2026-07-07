# fractals

A command-line tool that renders ASCII art fractals directly in your terminal.
It currently supports two fractal types: **Sierpinski triangles** (recursive
subdivision) and the **Mandelbrot set** (escape-time rendering with a
character gradient). Output size, recursion depth, iteration count, and the
character used for drawing are all configurable via flags.

## Installation

Requires Go 1.24+.

```bash
go install ./cmd/fractals
```

This builds the `fractals` binary and installs it to your `$GOPATH/bin` (or
`$GOBIN`). Make sure that directory is on your `PATH`.

Alternatively, build a local binary without installing:

```bash
go build -o fractals ./cmd/fractals
```

## Usage

```bash
fractals --help
fractals sierpinski --help
fractals mandelbrot --help
```

### `sierpinski`

Renders a Sierpinski triangle using recursive subdivision.

| Flag       | Default | Description                                   |
|------------|---------|------------------------------------------------|
| `--size`   | `32`    | Width of the triangle base in characters       |
| `--depth`  | `5`     | Recursion depth                                |
| `--char`   | `*`     | Character to use for filled points (exactly 1 char) |

```bash
fractals sierpinski --size 16 --depth 3
```

```
               *
              ***
             *   *
            *** ***
           *       *
          ***     ***
         *   *   *   *
        *** *** *** ***
       *               *
      ***             ***
     *   *           *   *
    *** ***         *** ***
   *       *       *       *
  ***     ***     ***     ***
 *   *   *   *   *   *   *   *
*** *** *** *** *** *** *** ***
```

Custom character:

```bash
fractals sierpinski --size 16 --depth 3 --char '#'
```

```
               #
              ###
             #   #
            ### ###
           #       #
          ###     ###
         #   #   #   #
        ### ### ### ###
       #               #
      ###             ###
     #   #           #   #
    ### ###         ### ###
   #       #       #       #
  ###     ###     ###     ###
 #   #   #   #   #   #   #   #
### ### ### ### ### ### ### ###
```

### `mandelbrot`

Renders the Mandelbrot set as ASCII art. By default it maps escape iteration
counts to a shading gradient (` .:-=+*#%@`, from empty to dense); pass
`--char` to instead draw a single character for points inside the set.

| Flag           | Default   | Description                                          |
|----------------|-----------|-------------------------------------------------------|
| `--width`      | `80`      | Output width in characters                            |
| `--height`     | `24`      | Output height in characters                            |
| `--iterations` | `100`     | Maximum iterations for escape calculation              |
| `--char`       | (gradient)| Single character to draw instead of the default gradient |

```bash
fractals mandelbrot --width 50 --height 16
```

```
                                 . @              
                                .=.               
                                @@@.              
                           .@-@@@@@@@+ ::         
                          .+@@@@@@@@@@@-          
                   .     .@@@@@@@@@@@@@@@.        
                  .@@@@@.@@@@@@@@@@@@@@@@:        
                +:@@@@@@@@@@@@@@@@@@@@@@.         
                +:@@@@@@@@@@@@@@@@@@@@@@.         
                  .@@@@@.@@@@@@@@@@@@@@@@:        
                   .     .@@@@@@@@@@@@@@@.        
                          .+@@@@@@@@@@@-          
                           .@-@@@@@@@+ ::         
                                @@@.              
                                .=.               
                                 . @              
```

With a single custom character:

```bash
fractals mandelbrot --width 40 --height 12 --char '#'
```

```
                                        
                          ##            
                         #####          
                     ###########        
               #    #############       
               #################        
               #################        
               #    #############       
                     ###########        
                         #####          
                          ##            
                                        
```

## Errors

Invalid flag values are rejected with a clear message on stderr and a
non-zero exit code:

```bash
$ fractals sierpinski --size -5
Error: size must be greater than 0, got -5

$ fractals sierpinski --char 'ab'
Error: char must be exactly one character, got "ab"
```
