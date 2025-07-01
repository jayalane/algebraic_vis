# Algebraic Numbers Visualizer

A visualization tool for algebraic numbers in the complex plane, showing the beautiful patterns formed by roots of integer polynomials.

## Original Author

This project is based on code originally written by [Stephen J. Brooks](https://en.wikipedia.org/wiki/User:Stephen_J._Brooks) and posted on Wikipedia. The original implementation demonstrated the fascinating distribution of algebraic numbers across the complex plane.

## What Are Algebraic Numbers?

Algebraic numbers are complex numbers that are roots of polynomial equations with integer coefficients. For example:
- √2 is algebraic (root of x² - 2 = 0)
- i is algebraic (root of x² + 1 = 0)
- ∛5 is algebraic (root of x³ - 5 = 0)

This visualizer generates polynomials of various "heights" (complexity measures) and plots their roots, colored by leading coefficient.

## Features

- **Interactive viewport**: Specify custom complex plane rectangles
- **Color coding**: Different colors for different leading coefficients
- **Scalable rendering**: Blob sizes adjust based on polynomial height
- **High performance**: Renders in under 1 second using parallel processing
- **Video animation**: Watch algebraic numbers fill the plane as complexity increases
- **Parallel processing**: Utilizes all CPU cores for fast computation

## Usage

```bash
make go                              # Build Go version
./algebraic_go                       # Default view (-2-2i to 2+2i)
./algebraic_go 0 -1 1 2              # Custom rectangle (0-i to 1+2i)
./algebraic_go --max-height 20       # Higher detail
./algebraic_go --help                # Show usage
```

### Video Animation

Generate animated videos showing how algebraic numbers progressively fill the complex plane:

```bash
# Create a video animation from height 2 to 12
./algebraic_go --video --max-height 12

# Faster animation with custom frame rate
./algebraic_go --video --fps 5 --max-height 8

# Custom viewport animation
./algebraic_go --video --max-height 10 -- -1 -1 1 1

# Custom output filename
./algebraic_go --video --output my_animation.mp4 --max-height 15
```

## Color Scheme

Colors indicate the **leading coefficient** of the polynomial (not the degree):

- **Red**: Leading coefficient 1 (algebraic integers)
- **Green**: Leading coefficient 2
- **Blue**: Leading coefficient 3
- **Yellow**: Leading coefficient 4
- **Magenta**: Leading coefficient 5
- **Cyan**: Leading coefficient 6
- **Orange**: Leading coefficient 7
- **Lime**: Leading coefficient 8
- **Hot Pink**: Leading coefficient 9
- **Purple**: Leading coefficient 10
- **White**: Higher coefficients

Point size decreases as the polynomial height (sum of coefficient magnitudes) increases.

## Mathematical Background

The visualization generates polynomials with integer coefficients up to a specified "height" (sum of absolute values of coefficients). For each polynomial, it finds all complex roots using Newton's method and plots them with:

- Size based on polynomial height (lower height = larger dots)
- Color based on leading coefficient
- Additive blending for overlapping points

## Requirements

- **Go**: 1.21+ (no external dependencies for static images)
- **Video generation**: ffmpeg (for animation mode)

## Output

Generates high-resolution PNG images or MP4 videos showing the intricate patterns formed by algebraic numbers in the complex plane, revealing the deep mathematical structure underlying polynomial equations.