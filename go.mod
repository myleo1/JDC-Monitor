module jd-mining-server

go 1.16

require (
	github.com/go-resty/resty/v2 v2.5.0
	github.com/mizuki1412/go-core-kit v0.0.0-20210312074446-ece47b1d8445
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/tidwall/gjson v1.6.8 // indirect
)

//replace github.com/mizuki1412/go-core-kit v0.0.0-20210224112604-6a0c16c64ed9 => ../mizuki/framework/core-kit
