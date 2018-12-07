package Collision

import "strconv"

func Collision(xstr string, ystr string, xblockstr string, yblockstr string, x2blockstr string) bool {
	x1, _ := strconv.ParseFloat(xstr, 64)
	y1, _ := strconv.ParseFloat(ystr, 64)
	x := int(x1)
	y := int(y1)
	xblock, _ := strconv.Atoi(xblockstr)
	x2block, _ := strconv.Atoi(x2blockstr)
	yblock, _ := strconv.Atoi(yblockstr)

	if (((x > xblock && x < xblock+60) || (x > xblock-200 && x < xblock+60-200) || (x > xblock-400 && x < xblock+60-400) || (x > xblock-600 && x < xblock+60-600)) && (y < yblock+60 && y > yblock)) {
		return true
	}
	if (((x > x2block && x < x2block+60) || (x > x2block+200 && x < x2block+60+200) || (x > x2block+400 && x < x2block+60+400) || (x > x2block+600 && x < x2block+60+600)) && (y < yblock+60+150 && y > yblock+150)) {
		return true
	}
	if (((x > xblock && x < xblock+60) || (x > xblock-200 && x < xblock+60-200) || (x > xblock-400 && x < xblock+60-400) || (x > xblock-600 && x < xblock+60-600)) && (y < yblock+60-250 && y > yblock-250)) {
		return true
	}

	return false
}
