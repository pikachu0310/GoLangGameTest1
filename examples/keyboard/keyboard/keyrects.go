// Copyright 2013 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gen.go using 'go generate'. DO NOT EDIT.

package keyboard

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var keyboardKeyRects = map[ebiten.Key]image.Rectangle{}

func init() {
	keyboardKeyRects[ebiten.KeyA] = image.Rect(48, 36, 64, 54)
	keyboardKeyRects[ebiten.KeyB] = image.Rect(120, 54, 136, 72)
	keyboardKeyRects[ebiten.KeyC] = image.Rect(88, 54, 104, 72)
	keyboardKeyRects[ebiten.KeyD] = image.Rect(80, 36, 96, 54)
	keyboardKeyRects[ebiten.KeyE] = image.Rect(72, 18, 88, 36)
	keyboardKeyRects[ebiten.KeyF] = image.Rect(96, 36, 112, 54)
	keyboardKeyRects[ebiten.KeyG] = image.Rect(112, 36, 128, 54)
	keyboardKeyRects[ebiten.KeyH] = image.Rect(128, 36, 144, 54)
	keyboardKeyRects[ebiten.KeyI] = image.Rect(152, 18, 168, 36)
	keyboardKeyRects[ebiten.KeyJ] = image.Rect(144, 36, 160, 54)
	keyboardKeyRects[ebiten.KeyK] = image.Rect(160, 36, 176, 54)
	keyboardKeyRects[ebiten.KeyL] = image.Rect(176, 36, 192, 54)
	keyboardKeyRects[ebiten.KeyM] = image.Rect(152, 54, 168, 72)
	keyboardKeyRects[ebiten.KeyN] = image.Rect(136, 54, 152, 72)
	keyboardKeyRects[ebiten.KeyO] = image.Rect(168, 18, 184, 36)
	keyboardKeyRects[ebiten.KeyP] = image.Rect(184, 18, 200, 36)
	keyboardKeyRects[ebiten.KeyQ] = image.Rect(40, 18, 56, 36)
	keyboardKeyRects[ebiten.KeyR] = image.Rect(88, 18, 104, 36)
	keyboardKeyRects[ebiten.KeyS] = image.Rect(64, 36, 80, 54)
	keyboardKeyRects[ebiten.KeyT] = image.Rect(104, 18, 120, 36)
	keyboardKeyRects[ebiten.KeyU] = image.Rect(136, 18, 152, 36)
	keyboardKeyRects[ebiten.KeyV] = image.Rect(104, 54, 120, 72)
	keyboardKeyRects[ebiten.KeyW] = image.Rect(56, 18, 72, 36)
	keyboardKeyRects[ebiten.KeyX] = image.Rect(72, 54, 88, 72)
	keyboardKeyRects[ebiten.KeyY] = image.Rect(120, 18, 136, 36)
	keyboardKeyRects[ebiten.KeyZ] = image.Rect(56, 54, 72, 72)
	keyboardKeyRects[ebiten.KeyArrowDown] = image.Rect(48, 126, 96, 144)
	keyboardKeyRects[ebiten.KeyArrowLeft] = image.Rect(0, 126, 48, 144)
	keyboardKeyRects[ebiten.KeyArrowRight] = image.Rect(96, 126, 144, 144)
	keyboardKeyRects[ebiten.KeyArrowUp] = image.Rect(48, 108, 96, 126)
	keyboardKeyRects[ebiten.KeyBackquote] = image.Rect(240, 0, 256, 18)
	keyboardKeyRects[ebiten.KeyBackslash] = image.Rect(224, 0, 240, 18)
	keyboardKeyRects[ebiten.KeyBackspace] = image.Rect(232, 18, 272, 36)
	keyboardKeyRects[ebiten.KeyBracketLeft] = image.Rect(200, 18, 216, 36)
	keyboardKeyRects[ebiten.KeyBracketRight] = image.Rect(216, 18, 232, 36)
	keyboardKeyRects[ebiten.KeyComma] = image.Rect(168, 54, 184, 72)
	keyboardKeyRects[ebiten.KeyDigit0] = image.Rect(176, 0, 192, 18)
	keyboardKeyRects[ebiten.KeyDigit1] = image.Rect(32, 0, 48, 18)
	keyboardKeyRects[ebiten.KeyDigit2] = image.Rect(48, 0, 64, 18)
	keyboardKeyRects[ebiten.KeyDigit3] = image.Rect(64, 0, 80, 18)
	keyboardKeyRects[ebiten.KeyDigit4] = image.Rect(80, 0, 96, 18)
	keyboardKeyRects[ebiten.KeyDigit5] = image.Rect(96, 0, 112, 18)
	keyboardKeyRects[ebiten.KeyDigit6] = image.Rect(112, 0, 128, 18)
	keyboardKeyRects[ebiten.KeyDigit7] = image.Rect(128, 0, 144, 18)
	keyboardKeyRects[ebiten.KeyDigit8] = image.Rect(144, 0, 160, 18)
	keyboardKeyRects[ebiten.KeyDigit9] = image.Rect(160, 0, 176, 18)
	keyboardKeyRects[ebiten.KeyEnter] = image.Rect(224, 36, 272, 54)
	keyboardKeyRects[ebiten.KeyEqual] = image.Rect(208, 0, 224, 18)
	keyboardKeyRects[ebiten.KeyEscape] = image.Rect(0, 0, 32, 18)
	keyboardKeyRects[ebiten.KeyMinus] = image.Rect(192, 0, 208, 18)
	keyboardKeyRects[ebiten.KeyPeriod] = image.Rect(184, 54, 200, 72)
	keyboardKeyRects[ebiten.KeyQuote] = image.Rect(208, 36, 224, 54)
	keyboardKeyRects[ebiten.KeySemicolon] = image.Rect(192, 36, 208, 54)
	keyboardKeyRects[ebiten.KeySlash] = image.Rect(200, 54, 216, 72)
	keyboardKeyRects[ebiten.KeySpace] = image.Rect(96, 72, 176, 90)
	keyboardKeyRects[ebiten.KeyTab] = image.Rect(0, 18, 40, 36)
	keyboardKeyRects[ebiten.KeyAlt] = image.Rect(64, 72, 96, 90)
	keyboardKeyRects[ebiten.KeyControl] = image.Rect(0, 36, 48, 54)
	keyboardKeyRects[ebiten.KeyShift] = image.Rect(0, 54, 56, 72)
}

func KeyRect(key ebiten.Key) (image.Rectangle, bool) {
	r, ok := keyboardKeyRects[key]
	return r, ok
}
