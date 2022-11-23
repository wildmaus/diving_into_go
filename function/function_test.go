package main

import (
	"math"
	"testing"
)

const MESSAGE = "Expect %[1]v to be equal %[2]v\n Expected: %[1]v\n Got: %[2]v\n"

var gameStampsTest []ScoreStamp

type testcases struct {
	offset        int
	expectedIndex int
}

/**
 * @notice func with same logic as getScore, but work with any gameStamps len
 * @return index of finded elem in gameStamps (instead of Score)
 */
func getIndex(gameStamps []ScoreStamp, offset int) int {
	switch {
	case offset <= 0 || gameStamps == nil || len(gameStamps) == 0:
		return 0
	case offset > gameStamps[len(gameStamps)-1].Offset:
		return len(gameStamps) - 1
	}
	start := int(math.Round(float64(offset) / 3.0))
	end := offset
	if end >= len(gameStamps) {
		end = len(gameStamps) - 1
	}
	for start < end {
		middle := (start + end) / 2
		switch {
		case gameStamps[middle].Offset == offset:
			return middle
		case gameStamps[middle].Offset < offset:
			start = middle + 1
		case gameStamps[middle].Offset > offset:
			end = middle - 1
		}
	}
	if gameStamps[start].Offset <= offset {
		return start
	}
	return start - 1
}

/**
 * @notice fun for generate gameStamps with exact step and length
 */
func generateGameStamps(length, step int) []ScoreStamp {
	stamps := make([]ScoreStamp, length)
	offset := 0
	for i := 0; i < length; i++ {
		stamps[i] = ScoreStamp{
			Offset: offset,
			Score: Score{
				Home: i,
				Away: 0,
			},
		}
		offset += step
	}
	return stamps
}

func TestMain(m *testing.M) {
	gameStampsTest = *fillScores()
	m.Run()
}

func TestGetScore(t *testing.T) {

	t.Run("Test with original gameStamps", func(t *testing.T) {
		t.Run("offset <= 0", func(t *testing.T) {
			offsets := []int{-10, 0}
			expectedScore := Score{}
			for _, offset := range offsets {
				score := getScore(gameStampsTest, offset)
				if score != expectedScore {
					t.Errorf(MESSAGE, expectedScore, score)
				}
			}
		})

		t.Run("Offset in the middle of gameStamps", func(t *testing.T) {
			offsets := []int{123, 2000, 15646}
			for _, offset := range offsets {
				resIndex := getIndex(gameStampsTest, offset)
				if gameStampsTest[resIndex].Offset > offset {
					t.Errorf("Expect offset %v to be less or equal %v", gameStampsTest[resIndex].Offset, offset)
				}
				if resIndex+1 < len(gameStampsTest) && gameStampsTest[resIndex+1].Offset <= offset {
					t.Errorf("Expect offset %v to be more %v", gameStampsTest[resIndex].Offset, offset)
				}
				expectedScore := gameStampsTest[resIndex].Score
				score := getScore(gameStampsTest, offset)
				if score != expectedScore {
					t.Errorf(MESSAGE, expectedScore, score)
				}
			}
		})

		t.Run("Offset >= right edge", func(t *testing.T) {
			rightOffset := gameStampsTest[TIMESTAMPS_COUNT-1].Offset
			offsets := []int{rightOffset, rightOffset + 1000}
			expectedScore := gameStampsTest[TIMESTAMPS_COUNT-1].Score
			for _, offset := range offsets {
				index := getIndex(gameStampsTest, offset)
				if index != TIMESTAMPS_COUNT-1 {
					t.Errorf(MESSAGE, TIMESTAMPS_COUNT-1, index)
				}
				score := getScore(gameStampsTest, offset)
				if score != expectedScore {
					t.Errorf(MESSAGE, expectedScore, score)
				}
			}
		})
	})

	t.Run("Positive with synthetic gameStamps", func(t *testing.T) {
		gameStemps := generateGameStamps(16, 3)

		t.Run("Existing offset", func(t *testing.T) {
			testcases := []testcases{
				{
					offset:        3,
					expectedIndex: 1,
				},
				{
					offset:        9,
					expectedIndex: 3,
				},
				{
					offset:        45,
					expectedIndex: 15,
				},
			}
			for _, tc := range testcases {
				index := getIndex(gameStemps, tc.offset)
				if index != tc.expectedIndex {
					t.Errorf(MESSAGE, tc.expectedIndex, index)
				}
			}
		})

		t.Run("Offset doesn't exist", func(t *testing.T) {
			testcases := []testcases{
				{
					offset:        2,
					expectedIndex: 0,
				},
				{
					offset:        4,
					expectedIndex: 1,
				},
				{
					offset:        44,
					expectedIndex: 14,
				},
			}
			for _, tc := range testcases {
				index := getIndex(gameStemps, tc.offset)
				if index != tc.expectedIndex {
					t.Errorf(MESSAGE, tc.expectedIndex, index)
				}
			}
		})
	})

	t.Run("Negative", func(t *testing.T) {
		t.Run("Nil as gameStamps", func(t *testing.T) {
			expectedScore := Score{}
			score := getScore(nil, 1)
			if score != expectedScore {
				t.Errorf(MESSAGE, expectedScore, score)
			}
		})

		t.Run("GameStamps with len < TIMESTAMPS_COUNT", func(t *testing.T) {
			lengths := []int{0, 1, TIMESTAMPS_COUNT - 1}
			expectedScore := Score{}
			for _, length := range lengths {
				gameStamps := generateGameStamps(length, 1)
				score := getScore(gameStamps, 4)
				if score != expectedScore {
					t.Errorf(MESSAGE, expectedScore, score)
				}
			}

		})
	})
}
