package phfonts

import (
	_ "embed"
)

var (
	//go:embed mplus-1p-regular.ttf
	MPlus1pRegular_ttf []byte

	//go:embed HackGen-Regular.ttf
	HackGenRegular_ttf []byte

	//go:embed NotoSansCJK-Regular.ttc
	NotoSansCJKRegular_ttc []byte

	//go:embed NotoSansCJK-Thin.ttc
	NotoSansCJKThin_ttc []byte
  //go:embed NotoSansMonoCJKjp-Regular.otf
  NotoSansMonoCJKjpRegular_otf []byte
)
