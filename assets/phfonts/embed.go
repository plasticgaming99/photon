package phfonts

import (
	_ "embed"
)

var (
	//go:embed mplus-1p-regular.ttf
	MPlus1pRegular_ttf []byte

	//go:embed HackGen-Regular.ttf
	HackGenRegular_ttf

	//go:embed NotoSansCJK-Regular.ttc
	NotoSansCJKRegular_ttc

	//go:embed NotoSansCJK-Thin.ttc
	NotoSansCJKThin_ttc
)
