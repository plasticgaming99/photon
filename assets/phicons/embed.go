package phicons

import (
	_ "embed"
)

var (
	//go:embed photon16.png
	PhotonIcon16 []byte
	//go:embed photon32.png
	PhotonIcon32 []byte
	//go:embed photon48.png
	PhotonIcon48 []byte
	//go:embed photon128.png
	PhotonIcon128 []byte
	//go:embed photonbig.png
	PhotonIconBig []byte
)
