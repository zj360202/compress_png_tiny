package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"os"
	"sort"
	"strconv"
	"strings"
	//"strconv"
	//"time"
	//"reflect"
)

//func Set(x, y int, c color.Color, p *image.Paletted) *image.Paletted {
//	//fmt.Println("1......")
//	if !(image.Point{x, y}.In(p.Rect)) {
//		return p
//	}
//	i := p.PixOffset(x, y)
//	p.Pix[i] = uint8(Index(c, p.Palette))
//	return p
//}
//
//func Index(c color.Color, p color.Palette) int {
//	// A batch version of this computation is in image/draw/draw.go.
//	var mod uint8 = 51
//	var splitNum = 255/mod + 1
//	cr, cg, cb, ca := c.RGBA()
//	cr8, cg8, cb8, ca8 := uint8(cr), uint8(cg), uint8(cb), uint8(ca)
//	ri, gi, bi, ai := cr8/mod, cg8/mod, cb8/mod, ca8/mod
//	rm, gm, bm, am := cr8%mod, cg8%mod, cb8%mod, ca8%mod
//	if rm > mod/2 {
//		ri += 1
//	}
//	if gm > mod/2 {
//		gi += 1
//	}
//	if bm > mod/2 {
//		bi += 1
//	}
//	if am > mod/2 {
//		ai += 1
//	}
//	//ret := int(ri * 36 + gi * 6 + bi)
//	ret := int(ri*splitNum*splitNum + gi*6 + bi)
//	return ret
//}
//func sqDiff(x, y uint32) uint32 {
//	d := x - y
//	return (d * d) >> 2
//}
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
func main() {
	filer, err := os.Open("1.png")
	if err != nil {
		fmt.Println(err)
	}
	defer filer.Close()
	subimg, err := png.Decode(filer)
	imgBounds := subimg.Bounds()
	newImg := image.NewPaletted(imgBounds, palette.Plan9)
	//pale := make(map[string]int)
	//maxIndex := 0
	newImg.Palette = make([]color.Color, 0)
	w, h := imgBounds.Max.X, imgBounds.Max.Y
	img := subimg.(*image.NRGBA)
	//img := subimg.(*image.RGBA)

	dataInfo := make(map[string]int)
	var modNum uint8 = 5
	for i := imgBounds.Min.X; i < w; i++ {
		for j := imgBounds.Min.Y; j < h; j++ {
			startIndex := (j*w + i) << 2
			pix1 := img.Pix[startIndex : startIndex+4]
			r1, b1, g1, a1 := pix1[0], pix1[1], pix1[2], pix1[3]
			ri, gi, bi, ai := r1/modNum, g1/modNum, b1/modNum, a1/modNum
			rm, gm, bm, am := r1%modNum, g1%modNum, b1%modNum, a1%modNum
			if rm > (modNum-1)>>2 {
				ri += 1
			}
			if gm > (modNum-1)>>2 {
				gi += 1
			}
			if bm > (modNum-1)>>2 {
				bi += 1
			}
			if am > (modNum-1)>>2 {
				ai += 1
			}
			rn := If(int(ri)*int(modNum) >= 256, 255, ri*modNum)
			bn := If(int(bi)*int(modNum) >= 256, 255, bi*modNum)
			gn := If(int(gi)*int(modNum) >= 256, 255, gi*modNum)
			an := If(int(ai)*int(modNum) >= 256, 255, ai*modNum)
			key := fmt.Sprintf("%03d", rn) + "," + fmt.Sprintf("%03d", bn) + "," + fmt.Sprintf("%03d", gn) + "," + fmt.Sprintf("%03d", an)
			if _, ok := dataInfo[key]; ok {
				dataInfo[key] += 1
			} else {
				dataInfo[key] = 1
			}
		}
	}

	//fmt.Println("dataInfo_len:", len(dataInfo))
	dataMaps := make(map[string]string)
	// 压缩到256个
	var dataKeys []string
	var dataValues []int
	for k, v := range dataInfo {
		dataKeys = append(dataKeys, k)
		dataValues = append(dataValues, v)
	}
	sort.Strings(dataKeys)
	//sort.Ints(dataValues)
	sort.Sort(sort.Reverse(sort.IntSlice(dataValues)))
	//fmt.Println("dataKeys:", dataKeys)
	thred := dataValues[255]
	repeatNum := 1
	// 判定重复数量
	for i := 254; i > 0; i-- {
		if dataValues[i] == thred {
			repeatNum += 1
		} else {
			break
		}
	}
	//fmt.Println("dataValues:", dataValues)
	//fmt.Println("thred:", thred)
	lastKey := ""
	var leaveKeys []string
	dataIds := make(map[string]int)
	ids, leaveRepeat := 0, repeatNum
	for _, dataK := range dataKeys {
		if dataInfo[dataK] >= thred {
			if leaveRepeat > 0 && dataInfo[dataK] == thred {
				leaveRepeat--
			} else if leaveRepeat == 0 && dataInfo[dataK] == thred {
				if lastKey != "" {
					dataMaps[dataK] = lastKey
				} else {
					leaveKeys = append(leaveKeys, dataK)
				}
				continue
			}
			if len(leaveKeys) != 0 {
				for _, lk := range leaveKeys {
					dataMaps[lk] = dataK
				}
			}
			vks := strings.Split(dataK, ",")
			//fmt.Println("vks", vks, vk, key, pix1, am, ai, an)
			rv, _ := strconv.Atoi(vks[0])
			bv, _ := strconv.Atoi(vks[1])
			gv, _ := strconv.Atoi(vks[2])
			av, _ := strconv.Atoi(vks[3])
			newImg.Palette = append(newImg.Palette, color.RGBA{uint8(rv), uint8(bv), uint8(gv), uint8(av)})
			dataMaps[dataK] = dataK
			dataIds[dataK] = ids
			//fmt.Println("ids:", ids, dataK, dataInfo[dataK], "[", rv, bv, gv, av, "]")
			ids += 1
			lastKey = dataK
			leaveKeys = nil
		} else {
			if lastKey != "" {
				dataMaps[dataK] = lastKey
			} else {
				leaveKeys = append(leaveKeys, dataK)
			}
		}
	}
	//fmt.Println("dataMaps:", len(newImg.Palette))
	for i := imgBounds.Min.X; i < w; i++ {
		for j := imgBounds.Min.Y; j < h; j++ {
			startIndex := (j*w + i) << 2
			pix1 := img.Pix[startIndex : startIndex+4]
			r1, b1, g1, a1 := pix1[0], pix1[1], pix1[2], pix1[3]
			ri, gi, bi, ai := r1/modNum, g1/modNum, b1/modNum, a1/modNum
			rm, gm, bm, am := r1%modNum, g1%modNum, b1%modNum, a1%modNum
			if rm > (modNum-1)>>2 {
				ri += 1
			}
			if gm > (modNum-1)>>2 {
				gi += 1
			}
			if bm > (modNum-1)>>2 {
				bi += 1
			}
			if am > (modNum-1)>>2 {
				ai += 1
			}

			rn := If(int(ri)*int(modNum) >= 256, 255, ri*modNum)
			bn := If(int(bi)*int(modNum) >= 256, 255, bi*modNum)
			gn := If(int(gi)*int(modNum) >= 256, 255, gi*modNum)
			an := If(int(ai)*int(modNum) >= 256, 255, ai*modNum)
			key := fmt.Sprintf("%03d", rn) + "," + fmt.Sprintf("%03d", bn) + "," + fmt.Sprintf("%03d", gn) + "," + fmt.Sprintf("%03d", an)
			vk, _ := dataMaps[key]
			vIndex := dataIds[vk]
			//if vIndex != 252{
			//	newImg.Pix[i] = 252
			//}else{
			//	newImg.Pix[i] = 252
			//}
			k := newImg.PixOffset(i, j)
			newImg.Pix[k] = uint8(vIndex)

			//if r1 == 255 || b1 == 255 || g1 == 255{
			//	fmt.Println("vIndex:", vIndex, "vk:", vk, "key:", key, "[", rn, bn, gn, an,"]","[", ri, bi, gi, ai,"]", vIndex, vk, newImg.Palette[vIndex])
			//	break
			//}
		}
		//if i == 0 {
		//	break
		//}
	}
	//fmt.Println(newImg.Palette)
	filew, err := os.Create("2.png")
	defer filew.Close()
	err = png.Encode(filew, newImg)
	//fmt.Println("client2:", (time.Now().UnixNano()/1e6 - start_time))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println()
}
