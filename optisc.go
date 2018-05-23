package dbscan

import "github.com/XanthusL/zset"


type Point interface {
	Distance(c interface{}) float64
	GetID() int64
}

func OPTISC(src []Point,minPts int, e float64) ([]Point,[]float64){
	coreDistance := make(map[int64]*zset.SortedSet)
	// calculate core distances
	for _, c := range src {
		coreDistance[c.GetID()] = getCoreDistance(c, src)
	}
	sorted := zset.New()
	result := make([]Point, 0)
	resultTable := make(map[int64]struct{})
	for len(src) != 0 {
		point := src[len(src)-1]
		src = src[:len(src)-1]
		_, score, _ := coreDistance[point.GetID()].GetDataByRank(int64(minPts), false)
		isCoreObj := score > 0 && score <= e
		if isCoreObj {
			sorted.Set(score, point.GetID(), point)
		}
		if sorted.Length() == 0 {
			continue
		}
		key, score, dat := sorted.GetDataByRank(0, false)
		sorted.Delete(key)
		if _, ok := resultTable[key]; !ok {
			resultTable[key] = struct{}{}
			result = append(result, dat.(Point))
		}
		length := coreDistance[key].Length()
		for i := int64(0); i < length; i++ {
			key2, score2, dat2 := coreDistance[key].GetDataByRank(i, false)
			if _, ok := resultTable[key2]; ok {
				continue
			}
			_, ok := sorted.GetData(key2)
			if ok {
				_, coreD, _ := sorted.GetRank(key2, false)
				if coreD > score2 {
					sorted.Set(score2, key2, dat2.(Point))
				}
			} else {
				sorted.Set(score2, key2, dat2.(Point))
			}
		}
	}
	dist:= make([]float64,len(result))
	for i,r:=range result{
		_, score, _ := coreDistance[r.GetID()].GetDataByRank(int64(minPts), false)
		dist[i] = score
	}
	return result,dist
}

func getCoreDistance(p Point, src []Point) *zset.SortedSet {
	set := zset.New()
	for _, target := range src {
		d := p.Distance(target)
		set.Set(d, target.GetID(), target)
	}
	return set
}
