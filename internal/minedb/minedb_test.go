package minedb

import (
	"math/rand"
	"minedb/mocks"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
)

func checkGet(t *testing.T, curDb *mineDb, key string, expectedValue any) {
	var valGet = expectedValue
	err := curDb.Get(key, &valGet)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	if !reflect.DeepEqual(valGet, expectedValue) {
		t.Errorf("With key = %s value %v is needed, but got: %v\n", key, expectedValue, valGet)
	}
}

func checkSet(t *testing.T, curDb *mineDb, key string, value any) {
	err := curDb.Set(key, value)
	if err != nil {
		t.Errorf("%s\n", err)
	}
}

type concurrencySafeMap struct {
	mu     sync.Mutex
	curMap map[string]any
}

func (curSafeMap *concurrencySafeMap) setVal(key string, val any) {
	curSafeMap.mu.Lock()
	defer curSafeMap.mu.Unlock()
	curSafeMap.curMap[key] = val
}

func (curSafeMap *concurrencySafeMap) getVal(key string) any {
	curSafeMap.mu.Lock()
	defer curSafeMap.mu.Unlock()
	return curSafeMap.curMap[key]
}

func (curSafeMap *concurrencySafeMap) getKeys() []string {
	curSafeMap.mu.Lock()
	defer curSafeMap.mu.Unlock()
	var curKeys []string = []string{}
	for k := range curSafeMap.curMap {
		curKeys = append(curKeys, k)
	}
	return curKeys
}

type concurrencySafeNum struct {
	mu     sync.Mutex
	curNum int
}

func (curSafeNum *concurrencySafeNum) addNum(val int) {
	curSafeNum.mu.Lock()
	defer curSafeNum.mu.Unlock()
	curSafeNum.curNum += val
}

func (curSafeNum *concurrencySafeNum) getNum() int {
	curSafeNum.mu.Lock()
	defer curSafeNum.mu.Unlock()
	return curSafeNum.curNum
}

func TestMineDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocks.NewMocklogger(ctrl)
	logger.EXPECT().Info("read database file to initialize data").MaxTimes(1)
	logger.EXPECT().Warn("could not read database file because it does not exist").MaxTimes(1)

	curLink, err := NewMineDb(logger)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	t.Run("Performance test", func(t *testing.T) {
		var testSize int = 1000
		for i := 0; i < testSize; i++ {
			go func() { // test for runtime errors
				curKey := strconv.Itoa(i)
				curVal := i
				checkSet(t, curLink, curKey, curVal)
			}()
		}
	})
	t.Run("Unit tests", func(t *testing.T) {
		var testSize int = 1000
		keys := []string{"Time", "Location", "Mode", "What", "Doctor", "Ostrich"}
		values := []any{float64(123.34), 38, "who", [3]int{1, 2, 3}, 342, "AAAA", nil} // these should probably be random but i dont care
		curRandSource := rand.NewSource(38)
		randFunc := rand.New(curRandSource)
		var correctMineDb = &(concurrencySafeMap{})
		correctMineDb.curMap = make(map[string]any)
		randomizedKeys := []string{}
		randomizedVals := []any{}
		var keyInd concurrencySafeNum
		var valInd concurrencySafeNum
		for i := 0; i < 10*testSize; i++ {
			randomizedKeys = append(randomizedKeys, keys[randFunc.Intn(len(keys))]) // Turns out that getting a randon number is not a concurrecny safe thing to do :(
			randomizedVals = append(randomizedVals, values[randFunc.Intn(len(values))])
		}
		for i := 0; i < testSize; i++ { // test without multiflow
			curKey := randomizedKeys[keyInd.getNum()]
			curVal := randomizedVals[valInd.getNum()]
			keyInd.addNum(1)
			valInd.addNum(1)
			correctMineDb.setVal(curKey, curVal)
			checkSet(t, curLink, curKey, curVal)
		}
		allKeys := correctMineDb.getKeys()
		for ind := range allKeys {
			var valToGet any = correctMineDb.getVal(allKeys[ind])
			checkGet(t, curLink, allKeys[ind], valToGet)
		}
	})
}
