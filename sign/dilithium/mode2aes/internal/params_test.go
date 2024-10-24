package internal

import (
	"testing"

	"github.com/JI-0/circl/sign/dilithium/internal/common"
)

// Tests specific to the current mode

func TestVectorDeriveUniformLeqEta(t *testing.T) {
	var p common.Poly
	var seed [64]byte
	p2 := common.Poly{
		8380416, 0, 2, 8380415, 1, 1, 0, 1, 0, 1, 8380415, 2,
		8380415, 8380415, 2, 2, 2, 1, 0, 2, 8380416, 1, 8380415,
		8380415, 8380416, 8380415, 8380416, 8380415, 1, 1, 0, 1,
		0, 1, 2, 8380416, 2, 1, 8380416, 1, 1, 2, 0, 8380416,
		8380416, 2, 0, 2, 8380415, 0, 1, 2, 1, 1, 1, 0, 8380415,
		1, 2, 8380415, 8380416, 1, 8380415, 0, 1, 8380416, 8380416,
		8380415, 0, 2, 8380415, 1, 8380416, 0, 8380416, 8380416,
		8380416, 2, 2, 1, 2, 8380415, 2, 0, 8380415, 8380415, 0,
		2, 8380415, 8380415, 1, 8380415, 2, 8380415, 0, 1, 2,
		8380415, 8380416, 8380415, 0, 8380416, 1, 0, 2, 0, 2,
		8380415, 8380416, 2, 1, 8380415, 1, 8380416, 1, 8380415,
		8380415, 0, 8380416, 0, 0, 0, 0, 0, 2, 1, 2, 0, 0, 8380415,
		8380416, 2, 0, 1, 8380416, 2, 1, 8380416, 2, 1, 8380416,
		0, 2, 8380416, 2, 0, 8380415, 0, 2, 0, 8380415, 1, 0,
		8380415, 2, 8380416, 8380416, 8380415, 0, 0, 8380416, 2,
		2, 1, 8380416, 2, 1, 2, 0, 8380415, 1, 0, 2, 2, 1, 0, 0,
		1, 2, 0, 2, 0, 2, 2, 0, 0, 2, 2, 8380416, 2, 2, 0, 8380415,
		1, 2, 2, 1, 1, 8380415, 8380415, 2, 2, 1, 8380416, 8380415,
		2, 1, 0, 8380416, 8380415, 8380415, 0, 1, 0, 8380416,
		8380416, 8380416, 8380416, 2, 8380415, 1, 8380415, 0, 1,
		0, 8380416, 2, 8380415, 2, 1, 2, 1, 1, 0, 8380415, 2,
		8380416, 8380416, 8380415, 8380415, 0, 2, 8380416, 1,
		8380416, 8380415, 8380416, 8380415, 2, 8380416, 2, 8380415,
		2, 2, 1, 8380415,
	}
	for i := 0; i < 64; i++ {
		seed[i] = byte(i)
	}
	PolyDeriveUniformLeqEta(&p, &seed, 30000)
	p.Normalize()
	if p != p2 {
		t.Fatalf("%v != %v", p, p2)
	}
}

func TestVectorDeriveUniformLeGamma1(t *testing.T) {
	var p, p2 common.Poly
	var seed [64]byte
	p2 = common.Poly{
		8340798, 8313384, 49077, 22486, 123481, 8288752, 36503,
		8340997, 8302174, 8258535, 13603, 17223, 8335009, 8345989,
		91340, 8349862, 83710, 72846, 89691, 8272088, 26276, 6832,
		29103, 8313655, 72326, 8321983, 113475, 95773, 17886,
		8365924, 8279095, 8343830, 8273195, 8310637, 8253500,
		8374711, 8262430, 17388, 8294137, 8262960, 8290290, 47349,
		44452, 128195, 8377719, 130632, 8256416, 8287230, 120060,
		8323823, 115401, 8351851, 97604, 47772, 8363419, 8353971,
		1956, 8267893, 8362705, 17686, 122170, 101229, 3317, 14205,
		8368014, 97101, 8360617, 111843, 8357331, 16215, 8346959,
		8313944, 8309613, 8348252, 64256, 8294208, 8318089, 8335255,
		8324894, 8273750, 27850, 8260308, 8258591, 80542, 8320495,
		36517, 8340794, 8304320, 8320157, 625, 8292418, 8317653,
		8275617, 8352781, 109921, 121642, 8291715, 129643, 111667,
		8325995, 8368715, 8283849, 8281348, 1417, 8336033, 100081,
		30984, 22277, 8307048, 55908, 50909, 8326533, 34891, 98542,
		121511, 2614, 125602, 59900, 120456, 8351260, 30124, 52065,
		124214, 48354, 8281081, 116665, 17218, 74568, 4798, 8274275,
		8328948, 8269139, 8908, 8276509, 8270063, 8370525, 8257669,
		36128, 8313115, 8325113, 8257382, 77895, 8288147, 8294769,
		8273027, 8370871, 57085, 59514, 82308, 71173, 6475, 8313311,
		27188, 35803, 8296637, 100553, 8333397, 19553, 8373991,
		8361935, 62433, 8300, 8371479, 8297954, 8352934, 8286,
		8355336, 8335507, 8370548, 8301039, 8270317, 26478, 113694,
		8296283, 8271234, 8250245, 8372668, 8284012, 8264500, 85893,
		8322354, 8358407, 130156, 52458, 8291251, 122476, 8308146,
		4252, 118400, 74123, 8333546, 70542, 8325370, 60510, 1874,
		8377673, 50805, 78992, 66936, 8266050, 8367830, 8342582,
		8268085, 65238, 61045, 8312728, 70547, 8309034, 5696, 118654,
		8330845, 29553, 68995, 70518, 8351920, 8269399, 128395,
		122804, 10848, 8291860, 8324935, 3842, 8265342, 8266117,
		8368377, 8311281, 24039, 8343875, 97893, 12670, 8370577,
		92482, 8288562, 8269568, 8371831, 8324316, 76758, 8327193,
		46615, 10323, 8337373, 101795, 88456, 25023, 8351520, 94650,
		8264851, 14881, 104171, 22607, 8379204, 8310533, 8324603,
		8299017, 128723, 8291421,
	}
	for i := 0; i < 64; i++ {
		seed[i] = byte(i)
	}
	PolyDeriveUniformLeGamma1(&p, &seed, 30000)
	p.Normalize()
	if p != p2 {
		t.Fatalf("%v != %v", p, p2)
	}
}
