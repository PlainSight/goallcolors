package main

import (
  "fmt"
  "math/rand"
	"image"
  "image/png"
  "os"
  "image/color"
)

const arraysize = 32

const (
  no childrenState = 0
  initialized childrenState = 1
  yes childrenState = 2
)

type childrenState byte

type superColor struct {
  location *octTree
  x int
  y int
  r int
  g int
  b int
}

func (color *superColor) getColorDistance (other *superColor) int {
  rd := color.r - other.r
  gd := color.g - other.g
  bd := color.b - other.b
  
  return rd*rd + gd*gd + bd*bd
}

type octTree struct {
  minx int
  miny int
  minz int
  maxx int
  maxy int
  maxz int
  parent *octTree
  colors [arraysize]*superColor
  children [8]*octTree
  size int
  hasChildren childrenState
}

func (tree *octTree) hasPointInside (color *superColor) bool {
  return tree.minx <= color.r && color.r < tree.maxx &&
   tree.miny <= color.g && color.g < tree.maxy &&
   tree.minz <= color.b && color.b < tree.maxz
}


func (tree *octTree) putColorInTree (color *superColor) {
  if tree.hasChildren == yes {
    tree.putColorInChildTree(color)
  } else {
    if tree.size < arraysize {
      tree.colors[tree.size] = color
      color.location = tree
    } else {
      tree.split()
      tree.putColorInChildTree(color)
    }
  }
  
  tree.size++
}

func (tree *octTree) putColorInChildTree (color *superColor) {
  for i := 0; i < 8; i++ {
    if tree.children[i].hasPointInside(color) {
      tree.children[i].putColorInTree(color)
      return
    }
  }
}

func (tree *octTree) removeFromTree (color *superColor) {
  for i := 0; i < arraysize; i++ {
    if tree.colors[i] == color {
      tree.colors[i] = tree.colors[tree.size - 1]
      break
    }
  }
  
  for t := tree; t != nil; t = t.parent {
    t.size--
    
    if t.size < arraysize/2 && t.hasChildren == 2 {
      tColorIndex := 0
      
      for i := 0; i < 8; i++ {
        for j := 0; j < t.children[i].size; j++ {
          tColorIndex++
          t.colors[tColorIndex] = t.children[i].colors[j]
        }
        if t.children[i].hasChildren == yes {
          t.children[i].hasChildren = initialized
        } else {
          t.children[i].hasChildren = no
        }
        t.children[i].size = 0
      }
      
      for i := 0; i < t.size; i++ {
        t.colors[i].location = t
      }
      
      t.hasChildren = initialized
    }
  }
}

func (tree *octTree) shouldVisit (nom *superColor, nearest *superColor) bool {
  mx := (tree.minx + tree.maxx) / 2
  my := (tree.miny + tree.maxy) / 2
  mz := (tree.minz + tree.maxz) / 2
  
  var cx, cy, cz int
  if nom.r < mx {
    cx = tree.minx
  } else {
    cx = tree.maxx
  }
  if nom.g < my {
    cy = tree.miny
  } else {
    cy = tree.maxy
  }
  if nom.b < mz {
    cz = tree.minz
  } else {
    cz = tree.maxz
  }
  
  var dx, dy, dz int
	if nom.r >= tree.maxx || nom.r < tree.minx {
    dx = nom.r - cx
  }
	if nom.g >= tree.maxy || nom.g < tree.miny {
	  dy = nom.g - cy 
  }
	if nom.b >= tree.maxz || nom.b < tree.minz {
	  dz = nom.b - cz    
  }

	distancesqr := dx*dx + dy*dy + dz*dz

	return nom.getColorDistance(nearest) > distancesqr;
}

func (tree *octTree) split() {
  midx := (tree.minx + tree.maxx) / 2
  midy := (tree.miny + tree.maxy) / 2
  midz := (tree.minz + tree.maxz) / 2
  
  if tree.hasChildren == no {
    tree.children[0] = &(octTree {
      minx: tree.minx, miny: tree.miny, minz: tree.minz,
      maxx: midx,      maxy: midy,      maxz: midz,
      parent: tree,
    })
    tree.children[1] = &(octTree {
      minx: tree.minx, miny: tree.miny, minz: midz,
      maxx: midx,      maxy: midy,      maxz: tree.maxz,
      parent: tree,
    })
    tree.children[2] = &(octTree {
      minx: tree.minx, miny: midy,      minz: tree.minz,
      maxx: midx,      maxy: tree.maxy, maxz: midz,
      parent: tree,
    })
    tree.children[3] = &(octTree {
      minx: tree.minx, miny: midy,      minz: midz,
      maxx: midx,      maxy: tree.maxy, maxz: tree.maxz,
      parent: tree,
    })
    tree.children[4] = &(octTree {
      minx: midx,      miny: tree.miny, minz: tree.minz,
      maxx: tree.maxx, maxy: midy,      maxz: midz,
      parent: tree,
    })
    tree.children[5] = &(octTree {
      minx: midx,      miny: tree.miny, minz: midz,
      maxx: tree.maxx, maxy: midy,      maxz: tree.maxz,
      parent: tree,
    })
    tree.children[6] = &(octTree {
      minx: midx,      miny: midy,      minz: tree.minz,
      maxx: tree.maxx, maxy: tree.maxy, maxz: midz,
      parent: tree,
    })
    tree.children[7] = &(octTree {
      minx: midx,      miny: midy,      minz: midz,
      maxx: tree.maxx, maxy: tree.maxy, maxz: tree.maxz,
      parent: tree,
    })
  }
  
  tree.hasChildren = yes
  
  for i := 0; i < arraysize; i++ {
    tree.putColorInChildTree(tree.colors[i])
  }
}

func (tree *octTree) findNearestColorInTree (nom *superColor, nearest *superColor) *superColor {
  if tree.size == 0 || nearest != nil && !tree.shouldVisit(nom, nearest) {
    return nearest
  }
  
  if tree.hasChildren != yes {
    for i := 0; i < tree.size; i++ {
      if nearest == nil || nom.getColorDistance(tree.colors[i]) < nom.getColorDistance(nearest) {
        nearest = tree.colors[i]
      }
    }
    return nearest
  }
  
  var bestChild int
  if 2 * nom.r > tree.minx + tree.maxx {
    bestChild += 4
  }
  if 2 * nom.g > tree.miny + tree.maxy {
    bestChild += 2
  }
  if 2 * nom.b > tree.minz + tree.maxz {
    bestChild ++
  }
  
  nearest = tree.children[bestChild].findNearestColorInTree(nom, nearest)
  
  for i := 0; i < 8; i++ {
    if i == bestChild {
      continue
    }
    nearest = tree.children[i].findNearestColorInTree(nom, nearest)
  }
  return nearest
}



func setPixel(img *image.RGBA, width int, height int, col *superColor, tree *octTree, r int) {
  set := false
  
  var openSpaces [8][2]int
  
  for !set {
    closestNeighbour := tree.findNearestColorInTree(col, nil)
    
    minx := closestNeighbour.x - 1
    maxx := closestNeighbour.x + 1
    miny := closestNeighbour.y - 1
    maxy := closestNeighbour.y + 1
    if minx < 0 { minx = 0 }
    if maxx >= width { maxx = width - 1 }
    if miny < 0 { miny = 0 }
    if maxy >= height { maxy = height - 1 }
    
    numOpen := 0
    
    for x := minx; x <= maxx; x++ {
      for y := miny; y <= maxy; y++ {
        a := img.RGBAAt(x, y).A
        if a != 255 {
          openSpaces[numOpen][0] = x
          openSpaces[numOpen][1] = y
          numOpen++
          set = true
        }
      }
    }
    
    if !set {
      closestNeighbour.location.removeFromTree(closestNeighbour)
    } else {
      if numOpen == 1 {
        closestNeighbour.location.removeFromTree(closestNeighbour)
      }
      
      placement := r % numOpen
      
      rgba := color.RGBA {
        R: uint8(col.r),
        G: uint8(col.g),
        B: uint8(col.b),
        A: 255,        
      }
      
      img.SetRGBA(openSpaces[placement][0], openSpaces[placement][1], rgba)
      
      col.x = openSpaces[placement][0]
      col.y = openSpaces[placement][1]
      
      tree.putColorInTree(col)
    }
  }
}

func main() {
  
  root := octTree {
    minx: 0,
    miny: 0,
    minz: 0,
    maxx: 256,
    maxy: 256,
    maxz: 256,
    parent: nil,
  }
  
  var colors [256*256*256]*superColor
  
  for i := 0; i < 256*256*256; i++ {
    c := superColor {
      r: i & 0x00FF0000 >> 16,
      g: i & 0x0000FF00 >> 8,
      b: i & 0x000000FF,
    }
    colors[i] = &c
  }
  
  width := 4096
  height := 4096
  
  todo := width * height
  
  fmt.Printf("Created colors...\n")
  
  rec := image.Rect(0, 0, width, height)
  img := image.NewRGBA(rec)
  
  firstX := width / 2
  firstY := height / 2
  
  firstColor := colors[0]
  
  rgba := color.RGBA {
    R: uint8(firstColor.r),
    G: uint8(firstColor.g),
    B: uint8(firstColor.b),
    A: 255,        
  }
  
  img.SetRGBA(firstX, firstY, rgba)
  
  firstColor.x = firstX
  firstColor.y = firstY
  
  root.putColorInTree(firstColor)
  
  for i := 1; i < todo; i++ {
    r := rand.Intn((256*256*256)-i)
    
    temp := colors[i]
    colors[i] = colors[r]
    colors[r] = temp
    
    setPixel(img, width, height, colors[i], &root, r)
  }
  
  f, _ := os.Create("allcolors.png")
  png.Encode(f, img)
  f.Close()
}