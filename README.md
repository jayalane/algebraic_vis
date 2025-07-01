# Algebraic Numbers Visualizer

A visualization tool for algebraic numbers in the complex plane, showing the beautiful patterns formed by roots of integer polynomials.

## Original Author

This project is based on code originally written by [Stephen J. Brooks](https://en.wikipedia.org/wiki/User:Stephen_J._Brooks) and posted on Wikipedia. The original implementation demonstrated the fascinating distribution of algebraic numbers across the complex plane.

## What Are Algebraic Numbers?

Algebraic numbers are complex numbers that are roots of polynomial equations with integer coefficients. For example:
- √2 is algebraic (root of x² - 2 = 0)
- i is algebraic (root of x² + 1 = 0)
- ∛5 is algebraic (root of x³ - 5 = 0)

This visualizer generates polynomials of various "heights" (complexity measures) and plots their roots, colored by polynomial degree.

## Features

- **Multiple implementations**: C (OpenGL), C (PNG), and Go versions
- **Interactive viewport**: Specify custom complex plane rectangles
- **Color coding**: Different colors for different polynomial degrees
- **Scalable rendering**: Blob sizes adjust based on polynomial height
- **High performance**: Go version renders in under 1 second

## Usage

### Go Version (Recommended)
```bash
make go                    # Build Go version
./algebraic_go            # Default view (-2-2i to 2+2i)
./algebraic_go 0 -1 1 2   # Custom rectangle (0-i to 1+2i)
./algebraic_go --help     # Show usage
```

### C Versions
```bash
make png                  # Build C PNG version
./algebraic_png           # Generates PPM file
./run.sh                  # Builds and converts to PNG

make                      # Build C OpenGL version
./run.sh opengl          # Interactive OpenGL display
```

## Color Scheme

- **Red**: Linear polynomials (degree 1)
- **Green**: Quadratic polynomials (degree 2)  
- **Blue**: Cubic polynomials (degree 3)
- **Yellow**: Quartic polynomials (degree 4)
- **Orange**: Quintic polynomials (degree 5)
- **Cyan**: Degree 6
- **Magenta**: Degree 7
- **Gray**: Degree 8
- **White**: Higher degrees

## Mathematical Background

The visualization generates polynomials with integer coefficients up to a specified "height" (sum of absolute values of coefficients). For each polynomial, it finds all complex roots using Newton's method and plots them with:

- Size based on polynomial height (lower height = larger dots)
- Color based on polynomial degree
- Additive blending for overlapping points

## Requirements

- **Go version**: Go 1.21+ (no external dependencies)
- **C versions**: GCC, OpenGL/GLUT, ImageMagick (for PNG conversion)

## Output

Generates high-resolution PNG images showing the intricate patterns formed by algebraic numbers in the complex plane, revealing the deep mathematical structure underlying polynomial equations.