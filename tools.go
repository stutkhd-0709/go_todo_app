//go:build tools

package main

// go runをすると、そのパッケージを持っていない場合はgo installして最新版を取得しに入ってしまう
// なので一旦go getしてgo.modにバージョン管理させておき、tools.goでそのパッケージを記述することでgo.modのバージョンのパッケージが使用されるという感じな気がする
import _ "github.com/matryer/moq"
