package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Point represents an algebraic number with its properties
type Point struct {
	Z              complex128 // The complex number
	H              int        // Height (complexity measure)
	O              int        // Order (degree of polynomial)
	LeadingCoeff   int        // Leading coefficient of the polynomial
}

// Config holds rendering parameters
type Config struct {
	Width, Height   int
	XMin, YMin      float64
	XMax, YMax      float64
	MaxHeight       int
	OutputFile      string
}

// findRootsInnerWithRand implements Newton's method for polynomial root finding with custom random source
func findRootsInnerWithRand(coeffs []complex128, order int, rng *rand.Rand) []complex128 {
	if order == 1 {
		if coeffs[1] != 0 {
			return []complex128{-coeffs[0] / coeffs[1]}
		}
		return nil
	}

	var roots []complex128
	const maxIters = 5000
	const tolerance = 1e-20

	// Start with a random initial guess
	root := complex(rng.Float64()*2-1, rng.Float64()*2-1)

	for iter := 0; iter < maxIters; iter++ {
		oldRoot := root

		// Compute f(root) and f'(root) using Horner's method
		f := coeffs[order]
		df := complex(0, 0)

		for i := order - 1; i >= 0; i-- {
			df = df*root + f
			f = f*root + coeffs[i]
		}

		if cmplx.Abs(df) < 1e-15 {
			// Derivative too small, try new starting point
			root = complex(rng.Float64()*2-1, rng.Float64()*2-1)
			continue
		}

		// Newton's method step
		root = root - f/df

		// Check convergence
		if cmplx.Abs(root-oldRoot) < tolerance {
			roots = append(roots, root)
			break
		}

		// Restart with new random point occasionally
		if iter%500 == 0 && iter > 0 {
			root = complex(rng.Float64()*2-1, rng.Float64()*2-1)
		}
	}

	// Deflate polynomial and find remaining roots
	if len(roots) > 0 {
		r := roots[0]
		// Synthetic division: reduce polynomial by factor (x - r)
		newCoeffs := make([]complex128, order)
		newCoeffs[order-1] = coeffs[order]
		for i := order - 2; i >= 0; i-- {
			newCoeffs[i] = coeffs[i+1] + r*newCoeffs[i+1]
		}

		// Find remaining roots
		remaining := findRootsInnerWithRand(newCoeffs, order-1, rng)
		roots = append(roots, remaining...)
	}

	return roots
}

// findRootsInner implements Newton's method for polynomial root finding (compatibility wrapper)
func findRootsInner(coeffs []complex128, order int) []complex128 {
	return findRootsInnerWithRand(coeffs, order, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// PolyWork represents work for processing a single polynomial
type PolyWork struct {
	coeffs       []complex128
	h            int
	order        int
	leadingCoeff int
}

// generateAlgebraicNumbers computes algebraic numbers up to given height using parallel processing
func generateAlgebraicNumbers(maxHeight int) []Point {
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d CPU cores for parallel computation\n", numWorkers)
	
	// Channel for work distribution
	workCh := make(chan PolyWork, 1000)
	resultCh := make(chan []Point, 1000)
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each worker gets its own random source for thread safety
			localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			
			for work := range workCh {
				// Process this polynomial
				roots := findRootsInnerWithRand(work.coeffs, work.order, localRand)
				
				var workPoints []Point
				for _, root := range roots {
					workPoints = append(workPoints, Point{
						Z:            root,
						H:            work.h,
						O:            work.order,
						LeadingCoeff: work.leadingCoeff,
					})
				}
				resultCh <- workPoints
			}
		}()
	}
	
	// Generate work items
	go func() {
		defer close(workCh)
		
		for h := 2; h <= maxHeight; h++ {
			if maxHeight > 15 {
				fmt.Printf("Processing height %d/%d...\n", h, maxHeight)
			}
			// Generate all possible coefficient patterns for height h
			for i := (1 << (h - 1)) - 1; i >= 0; i -= 2 { // Step by 2 to avoid leading coefficient 0
				// Convert bit pattern to coefficient magnitudes
				coeffMags := make([]int, h)
				k := 0
				for j := h - 2; j >= 0; j-- {
					if (i>>j)&1 == 1 {
						coeffMags[k]++
					} else {
						k++
						if k < h {
							coeffMags[k] = 0
						}
					}
				}

				if k == 0 {
					continue // Invalid polynomial
				}

				order := k

				// Count non-zero coefficients for sign combinations
				nonZero := 0
				for j := 0; j <= order; j++ {
					if coeffMags[j] != 0 {
						nonZero++
					}
				}

				// Generate all sign combinations
				for signs := 0; signs < (1 << (nonZero - 1)); signs++ {
					// Build coefficient array
					coeffs := make([]complex128, order+1)
					signBit := 0

					for j := 0; j <= order; j++ {
						if coeffMags[j] == 0 || j == order {
							coeffs[j] = complex(float64(coeffMags[j]), 0)
						} else {
							sign := 1
							if (signs>>signBit)&1 == 1 {
								sign = -1
							}
							coeffs[j] = complex(float64(sign*coeffMags[j]), 0)
							signBit++
						}
					}

					// Send work to channel
					workCh <- PolyWork{
						coeffs:       coeffs,
						h:            h,
						order:        order,
						leadingCoeff: coeffMags[order],
					}
				}
			}
		}
	}()
	
	// Collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	
	var allPoints []Point
	stats := struct{ eqns, roots int }{}
	
	for points := range resultCh {
		stats.eqns++
		stats.roots += len(points)
		allPoints = append(allPoints, points...)
	}

	fmt.Printf("Generated: eqns=%d roots=%d\n", stats.eqns, stats.roots)
	return allPoints
}

// drawBlob draws a gaussian blob at the specified location with proper falloff
func drawBlob(img *image.RGBA, x, y int, radius float64, col color.RGBA) {
	bounds := img.Bounds()
	r := int(radius + 5) // Extend more for larger blobs

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			px, py := x+dx, y+dy
			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}

			dist := math.Sqrt(float64(dx*dx + dy*dy))
			// Gaussian falloff with wider spread for dramatic glow effect
			sigma := radius / 2.5 // Wider gaussian
			intensity := math.Exp(-dist*dist / (2 * sigma * sigma))
			
			if intensity > 0.005 { // Lower threshold for more glow
				// Get existing pixel
				existing := img.RGBAAt(px, py)
				
				// Add new color with intensity (additive blending)
				newR := float64(existing.R) + float64(col.R)*intensity
				newG := float64(existing.G) + float64(col.G)*intensity
				newB := float64(existing.B) + float64(col.B)*intensity
				
				// Clamp to 255
				if newR > 255 { newR = 255 }
				if newG > 255 { newG = 255 }
				if newB > 255 { newB = 255 }
				
				img.SetRGBA(px, py, color.RGBA{
					R: uint8(newR),
					G: uint8(newG), 
					B: uint8(newB),
					A: 255,
				})
			}
		}
	}
}

// getColorForLeadingCoeff returns color based on leading coefficient
// Red = 1 (algebraic integers), Green = 2, Blue = 3, Yellow = 4, etc.
func getColorForLeadingCoeff(coeff int) color.RGBA {
	switch coeff {
	case 1: return color.RGBA{255, 0, 0, 255}     // Red (algebraic integers)
	case 2: return color.RGBA{0, 255, 0, 255}     // Green
	case 3: return color.RGBA{0, 0, 255, 255}     // Blue
	case 4: return color.RGBA{255, 255, 0, 255}   // Yellow
	case 5: return color.RGBA{255, 0, 255, 255}   // Magenta
	case 6: return color.RGBA{0, 255, 255, 255}   // Cyan
	case 7: return color.RGBA{255, 128, 0, 255}   // Orange
	case 8: return color.RGBA{128, 255, 0, 255}   // Lime
	case 9: return color.RGBA{255, 0, 128, 255}   // Hot pink
	case 10: return color.RGBA{128, 0, 255, 255}  // Purple
	default: return color.RGBA{255, 255, 255, 255} // White for higher coefficients
	}
}

func renderImage(points []Point, config Config) error {
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	
	// Fill background with black
	bounds := img.Bounds()
	black := color.RGBA{0, 0, 0, 255}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, black)
		}
	}
	
	xRange := config.XMax - config.XMin
	yRange := config.YMax - config.YMin
	
	fmt.Printf("Rendering %d points to %dx%d image...\n", len(points), config.Width, config.Height)
	
	for _, point := range points {
		// Skip points outside viewport
		x, y := real(point.Z), imag(point.Z)
		if x < config.XMin || x > config.XMax || y < config.YMin || y > config.YMax {
			continue
		}
		
		// Transform to screen coordinates
		screenX := int((x - config.XMin) / xRange * float64(config.Width))
		screenY := int((config.YMax - y) / yRange * float64(config.Height)) // Flip Y
		
		// Calculate blob size (lower height = larger dots) - much larger like Wikipedia image
		k1 := 25.0 * (4.0 / xRange) // Much larger base size
		k2 := 0.5
		radius := k1 * math.Pow(k2, float64(point.H-3))
		
		if radius < 3.0 { radius = 3.0 }  // Larger minimum
		if radius > 80 { radius = 80 }    // Much larger maximum
		
		// Color based on leading coefficient (not degree!)
		color := getColorForLeadingCoeff(point.LeadingCoeff)
		drawBlob(img, screenX, screenY, radius, color)
	}
	
	// Save as PNG
	file, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}
	
	fmt.Printf("Saved image to %s\n", config.OutputFile)
	return nil
}

func printUsage(progName string) {
	fmt.Printf("Usage: %s [flags] [x_min y_min x_max y_max]\n", progName)
	fmt.Printf("  Renders algebraic numbers in the complex plane rectangle from (x_min + y_min*i) to (x_max + y_max*i)\n")
	fmt.Printf("\nFlags:\n")
	fmt.Printf("  --max-height N    Maximum polynomial height (complexity). Higher = more detail but slower (default: 15)\n")
	fmt.Printf("  --help, -h        Show this help message\n")
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s                           # Default view (-2-2i to 2+2i), height 15\n", progName)
	fmt.Printf("  %s --max-height 20           # Higher detail\n", progName)
	fmt.Printf("  %s 0 -1 1 2                  # Custom rectangle (0-i to 1+2i)\n", progName)
	fmt.Printf("  %s --max-height 25 -1 -1 1 1 # High detail, zoomed view\n", progName)
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Define flags
	maxHeight := flag.Int("max-height", 15, "Maximum polynomial height (complexity). Higher = more detail but slower")
	help := flag.Bool("h", false, "Show help message")
	helpLong := flag.Bool("help", false, "Show help message")
	
	// Custom usage function
	flag.Usage = func() {
		printUsage(os.Args[0])
	}
	
	// Parse flags
	flag.Parse()
	
	if *help || *helpLong {
		printUsage(os.Args[0])
		return
	}
	
	config := Config{
		Width:      1200,
		Height:     800,
		XMin:       -2.0,
		YMin:       -2.0,
		XMax:       2.0,
		YMax:       2.0,
		MaxHeight:  *maxHeight,
		OutputFile: "algebraic_numbers.png",
	}
	
	// Parse remaining positional arguments for viewport
	args := flag.Args()
	
	if len(args) == 4 {
		var err error
		if config.XMin, err = strconv.ParseFloat(args[0], 64); err != nil {
			log.Fatalf("Invalid x_min: %v", err)
		}
		if config.YMin, err = strconv.ParseFloat(args[1], 64); err != nil {
			log.Fatalf("Invalid y_min: %v", err)
		}
		if config.XMax, err = strconv.ParseFloat(args[2], 64); err != nil {
			log.Fatalf("Invalid x_max: %v", err)
		}
		if config.YMax, err = strconv.ParseFloat(args[3], 64); err != nil {
			log.Fatalf("Invalid y_max: %v", err)
		}
		
		if config.XMin >= config.XMax || config.YMin >= config.YMax {
			log.Fatal("Error: Invalid rectangle. x_min must be < x_max and y_min must be < y_max")
		}
	} else if len(args) != 0 {
		fmt.Println("Error: Wrong number of positional arguments")
		printUsage(os.Args[0])
		os.Exit(1)
	}
	
	// Validate max height
	if *maxHeight < 2 {
		log.Fatal("Error: max-height must be at least 2")
	}
	if *maxHeight > 30 {
		fmt.Printf("Warning: max-height %d is very high and may take a long time\n", *maxHeight)
	}
	
	fmt.Printf("Rendering complex plane from (%.2f + %.2fi) to (%.2f + %.2fi)\n",
		config.XMin, config.YMin, config.XMax, config.YMax)
	
	fmt.Println("Calculating algebraic numbers...")
	points := generateAlgebraicNumbers(config.MaxHeight)
	
	if err := renderImage(points, config); err != nil {
		log.Fatalf("Failed to render image: %v", err)
	}
}
