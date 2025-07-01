# Claude Development Context

## Project Overview
This is an algebraic numbers visualizer that generates beautiful visualizations of complex numbers that are roots of integer polynomials. The project started with broken C code from Wikipedia and has been fixed and enhanced with multiple implementations.

## Current State

### Working Implementations
1. **Go version** (`algebraic.go`) - **RECOMMENDED**
   - Fast, clean implementation (~0.9s runtime)
   - Direct PNG output via Go's image library
   - Native complex number support
   - Command line arguments for custom viewports
   - Black background for better aesthetics

2. **C PNG version** (`png_version.c`)
   - Fixed from original broken Wikipedia code
   - Outputs PPM format (converts to PNG via ImageMagick)
   - Same command line interface as Go version

3. **C OpenGL version** (`c.c`)
   - Interactive real-time visualization
   - Fixed OpenGL/GLUT implementation
   - Has some rendering issues (black screen problems)

### Build System
- `Makefile` supports all versions: `make go`, `make png`, `make all`
- `run.sh` script handles building and running with format conversion

## Key Technical Details

### Algorithm
1. **Polynomial Generation**: Creates integer polynomials up to specified height (complexity measure)
2. **Root Finding**: Uses Newton's method with deflation for multiple roots
3. **Rendering**: 
   - Colors by polynomial degree (red=linear, green=quadratic, etc.)
   - Blob size inversely proportional to polynomial height (simpler = larger)
   - Additive blending for overlapping points
   - Black background for aesthetic appeal

### Command Line Interface
```bash
./algebraic_go [x_min y_min x_max y_max]
```
- Default: complex plane from (-2-2i) to (2+2i)
- Example: `./algebraic_go 0 -1 1 2` shows rectangle from 0-i to 1+2i

### Performance
- Go version: ~0.9 seconds, generates ~35K points
- Direct PNG output, no external dependencies
- 1200x800 resolution by default

## Development History
1. Started with broken Wikipedia C code (missing headers, Windows-specific)
2. Fixed C compilation issues for Linux
3. Debugged OpenGL rendering problems
4. Created PNG-only C version to avoid display issues  
5. Rewrote in Go for better performance and maintainability
6. Added command line arguments and black background

## Future Enhancements
- Web version with interactive zoom/pan
- Higher polynomial heights for more detail
- Color customization options
- Animation/zooming capabilities
- Mathematical analysis tools (density plots, etc.)

## Files to Preserve
- `algebraic.go` - Main Go implementation
- `go.mod` - Go module file
- `Makefile` - Build system
- `README.md` - User documentation
- Original C files are kept for reference but Go version is preferred

## Key Dependencies
- Go 1.21+ (standard library only)
- For C versions: GCC, OpenGL/GLUT, ImageMagick

## Testing Commands
```bash
make go && ./algebraic_go                    # Test default view
./algebraic_go 0 -1 1 2                     # Test custom viewport
./algebraic_go --help                       # Test help
time ./algebraic_go                         # Performance test
```