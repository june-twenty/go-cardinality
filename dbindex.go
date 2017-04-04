package main

import (
	"github.com/hidai620/go-mysql-study/config"
	. "github.com/hidai620/go-mysql-study/dbindex"
	"github.com/hidai620/go-mysql-study/option"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)

	// コマンドラインオプションのパース
	opt, err := option.Parse()
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("option: %#v", opt)

	//　設定ファイルの読み込み
	conf, err := config.Load(opt.ConfigPath)
	if err != nil {
		logger.Println(err)
		return
	}

	// DB接続
	db, err := Connect(conf)
	if err != nil {
		logger.Println(err)
		return
	}
	defer db.Close()

	// 管理スキーマの取得
	informationSchema := NewInformationSchema(db)

	// テーブル単位の件数の取得
	tableRows, err := informationSchema.TableRows(conf.Database, opt.TableNames)
	if err != nil {
		logger.Println(err)
		return
	}

	if len(tableRows) != 0 {
		// カラムの取得
		columns, err := informationSchema.TableColumns(conf.Database, opt.TableNames)
		if err != nil {
			logger.Println(err)
			return
		}

		// 出力先の設定
		writer := getWriter(opt.Out, conf)
		err = writer.WriteDDL(columns, tableRows)
		if err != nil {
			logger.Println(err)
			return
		}
	}
}

// getWriter returns Writer according to command line argument.
func getWriter(out option.OutputType, config *config.Config) Writer {
	switch out {
	case option.CONSOLE:
		return NewConsole(os.Stdout, config)
	case option.CSV:
		return NewCSV(os.Stdout, config)
	default:
		return NewConsole(os.Stdout, config)
	}
}
