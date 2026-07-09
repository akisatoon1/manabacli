// manabaライブラリを本番環境で使うための実装.
//
// manabaは外部と通信する必要があるので, テストがしづらい.
// そのためmanabaのLoginとUploadFileをInterfaceにして, 依存を切り出した.
// しかし, manaba package自体がInterfaceを満たすことはできないので,
// Interfaceを満たすstruct型をここで定義する.
package manabaclient

import (
	"net/http/cookiejar"

	"github.com/akisatoon1/manaba"
)

type Client struct{}

func (c Client) Login(jar *cookiejar.Jar, username string, password string) error {
	return manaba.Login(jar, username, password)
}

func (c Client) UploadFile(jar *cookiejar.Jar, url string, filePath string) error {
	return manaba.UploadFile(jar, url, filePath)
}
