package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type cell struct {
	node    int
	element *canvas.Rectangle
	queued  bool
	wall    bool
	top     *canvas.Rectangle
	left    *canvas.Rectangle
	right   *canvas.Rectangle
	bottom  *canvas.Rectangle
	start   bool
	target  bool
	visited bool
	prior   int
}

var start_cell int = 32
var end_cell int = 161
var numberOfCol int = 15

var cellTable []cell
var queue []cell

// color variable

var wall_color color.NRGBA = color.NRGBA{R: 0, G: 100, B: 100, A: 255}
var visited_color color.NRGBA = (color.NRGBA{R: 30, G: 144, B: 255, A: 255})
var queue_color color.NRGBA = (color.NRGBA{R: 0, G: 0, B: 139, A: 255})
var red color.NRGBA = color.NRGBA{R: 32, G: 178, B: 170, A: 255}
var blue color.NRGBA = color.NRGBA{R: 0, G: 255, B: 255, A: 255}

// searching state
var searching bool = false

func main() {
	a := app.New()
	w := a.NewWindow("Djikstra algorithm")

	// Grid creation
	grid := createGrid()
	// Initialization of the start and end cell
	grid.Objects[start_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	grid.Objects[end_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	grid.Objects[start_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).Refresh()
	grid.Objects[end_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).Refresh()

	// Start button under the grid

	start := widget.NewButton("Start", func() {
		setNeighbours(grid)
		playGame(grid)

	}) // button widget

	// Window setup

	vBox := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), grid, start)
	start.Resize(fyne.NewSize(32, 32))
	w.SetContent(vBox)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()

}

func createGrid() *fyne.Container {
	grid := container.NewGridWithColumns(numberOfCol)
	counter := 0

	// Here we create the cells grid and add it to the container
	for y := 0; y < numberOfCol; y++ {
		for x := 0; x < numberOfCol; x++ {
			btn_color := canvas.NewRectangle(
				color.NRGBA{R: 255, G: 255, B: 240, A: 255})

			btn := widget.NewButton("", nil) // button widget
			btn.OnTapped = func() {
				if btn_color.FillColor == wall_color {
					btn_color.FillColor = color.NRGBA{R: 255, G: 255, B: 240, A: 255}
					btn_color.Refresh()
				} else if btn_color.FillColor == (color.NRGBA{R: 255, G: 255, B: 240, A: 255}) {
					btn_color.FillColor = wall_color
					btn_color.Refresh()
				}

			}

			container1 := container.New(
				layout.NewMaxLayout(),
				btn_color,
				btn,
			)
			grid.Add(container1)

			c := cell{
				node:    counter,
				element: grid.Objects[counter].(*fyne.Container).Objects[0].(*canvas.Rectangle),
			}

			cellTable = append(cellTable, c)

			counter++
		}
	}
	return grid
}

func setNeighbours(grid *fyne.Container) {

	//Here we define all neighbours of each cell

	for i := range cellTable {

		if cellTable[i].element.FillColor == wall_color {
			cellTable[i].wall = true
		}

		if i > numberOfCol {
			cellTable[i].top = grid.Objects[i-15].(*fyne.Container).Objects[0].(*canvas.Rectangle)
		} else {
			cellTable[i].top = nil
		}
		if i < (numberOfCol*numberOfCol - numberOfCol) {
			cellTable[i].bottom = grid.Objects[i+15].(*fyne.Container).Objects[0].(*canvas.Rectangle)
		} else {
			cellTable[i].bottom = nil
		}
		if i%numberOfCol != 0 && i < 224 {
			cellTable[i].right = grid.Objects[i+1].(*fyne.Container).Objects[0].(*canvas.Rectangle)
		} else {
			cellTable[i].right = nil
		}
		if i%numberOfCol != 0 && i > 0 {
			cellTable[i].left = grid.Objects[i-1].(*fyne.Container).Objects[0].(*canvas.Rectangle)
		} else {
			cellTable[i].left = nil
		}

	}

	// Here we apply the start and target settings for the starter cell
	cellTable[end_cell].target = true
	cellTable[start_cell].start = true
	cellTable[start_cell].visited = true
	queue = append(queue, cellTable[start_cell])
}

func playGame(grid *fyne.Container) {

	running := true
	searching = true
	i := 1
	for running {
		// While queue is not empty
		if len(queue) > 0 && searching {
			current_cell := queue[0]
			queue = queue[1:]
			current_cell.visited = true

			// If cell is target cell
			if cellTable[current_cell.node].target {
				searching = false
				for current_cell.node != start_cell {
					cellTable[current_cell.prior].element.FillColor = (color.NRGBA{R: 255, G: 255, B: 0, A: 255})
					cellTable[current_cell.prior].element.Refresh()
					current_cell = cellTable[current_cell.prior]
				}
				running = false
			} else {

				// If cell is not target cell

				// We give the visited color to the cell
				current_cell.element.FillColor = visited_color
				current_cell.element.Refresh()

				// TOP NEIGHBOUR
				top := current_cell.top
				if top != nil {

					isQueued := cellTable[current_cell.node-15].queued
					isWall := cellTable[current_cell.node-15].wall

					// Here we check if the cell is not a wall and not queued
					if !isQueued && !isWall {
						// We add the cell to the queue
						top := &cellTable[current_cell.node-15]
						top.prior = current_cell.node
						top.queued = true

						if !top.target {
							top.element.FillColor = queue_color
							canvas.Refresh(top.element)
						}

						queue = append(queue, *top)
					}
				}

				// RIGHT NEIGHBOUR

				right := current_cell.right
				if right != nil {

					isQueued := cellTable[current_cell.node+1].queued
					isWall := cellTable[current_cell.node+1].wall

					// Here we check if the cell is not a wall and not queued
					if !isQueued && !isWall && current_cell.node%15 != 14 {
						// We add the cell to the queue
						right := &cellTable[current_cell.node+1]
						right.prior = current_cell.node
						right.queued = true
						if !right.target {
							right.element.FillColor = queue_color
							canvas.Refresh(right.element)
						}

						queue = append(queue, *right)
					}
				}

				// LEFT NEIGHBOUR

				left := current_cell.left
				if left != nil {

					isQueued := cellTable[current_cell.node-1].queued
					isWall := cellTable[current_cell.node-1].wall

					// Here we check if the cell is not a wall and not queued
					if !isQueued && !isWall {
						// We add the cell to the queue
						left := &cellTable[current_cell.node-1]
						left.prior = current_cell.node
						left.queued = true

						if !left.target {
							left.element.FillColor = queue_color
							canvas.Refresh(left.element)
						}

						queue = append(queue, *left)
					}
				}

				// BoTTOM NEIGHBOUR

				bottom := current_cell.bottom
				if bottom != nil {

					isQueued := cellTable[current_cell.node+15].queued
					isWall := cellTable[current_cell.node+15].wall

					// Here we check if the cell is not a wall and not queued
					if !isQueued && !isWall && current_cell.node < 224 {
						// We add the cell to the queue
						bottom := &cellTable[current_cell.node+15]
						bottom.prior = current_cell.node
						bottom.queued = true
						if !bottom.target {
							bottom.element.FillColor = queue_color
							canvas.Refresh(bottom.element)

						}

						queue = append(queue, *bottom)
					}
				}

			}
		}

		// Reset the color of the starter cell and the target cell
		grid.Objects[start_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
		grid.Objects[end_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
		grid.Objects[start_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).Refresh()
		grid.Objects[end_cell].(*fyne.Container).Objects[0].(*canvas.Rectangle).Refresh()

		// Add a 500microsecond delay to make the display cool
		time.Sleep(500 * time.Microsecond)
		i++
	}

}
