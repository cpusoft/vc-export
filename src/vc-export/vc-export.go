package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/cpusoft/goutil/conf"
	iputil "github.com/cpusoft/goutil/iputil"
	jsonutil "github.com/cpusoft/goutil/jsonutil"
	_ "github.com/cpusoft/goutil/logs"
	osutil "github.com/cpusoft/goutil/osutil"
	xormdb "github.com/cpusoft/goutil/xormdb"
)

type Export struct {
	Asn       int    `json:"asn" xorm:"asn"`
	Prefix    string `json:"prefix"`
	MaxLength int    `json:"maxLength"   xorm:"maxLength"`
	Ta        string `json:"ta"`

	AddressPrefix []byte `json:"-" xorm:"addressPrefix"`
	PrefixLength  int    `json:"-"  xorm:"prefixLength"`
	DirName       string `json:"-"  xorm:"dirName"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage:\r\n1. get json: ./vc-export ./file1.json\r\n2. get csv: ./vc-export ./file1.csv")
		return
	}

	file := os.Args[1]

	// start mysql
	err := xormdb.InitMySql()
	if err != nil {
		fmt.Println("InitMySql failed, err is ", err)
		return
	}
	defer xormdb.XormEngine.Close()

	talMap := make(map[string]string)
	talMap["ca.rg.net"] = "ca.rg.net"
	talMap["repository.lacnic.net"] = "repository.lacnic.net"
	talMap["rpki.afrinic.net"] = "rpki.afrinic.net"
	talMap["rpki.apnic.net"] = "rpki.apnic.net"
	talMap["rpki.arin.net"] = "rpki.arin.net"
	talMap["rpkica.twnic.tw"] = "rpkica.twnic.tw"
	talMap["rpki.cnnic.cn"] = "rpki.cnnic.cn"
	talMap["rpki-repository.nic.ad.jp"] = "rpki-repository.nic.ad.jp"
	talMap["rpki.ripe.net"] = "rpki.ripe.net"

	exports := make([]Export, 0)
	sql := `SELECT rpki_roa.asn as asn, 
		(rpki_roa_prefix.prefix) as addressPrefix, 
		rpki_roa_prefix.prefix_length as prefixLength, 
		rpki_roa_prefix.prefix_max_length as maxLength,
		rpki_dir.dirname 
		from rpki_roa, rpki_roa_prefix , rpki_dir 
		where rpki_roa.local_id = rpki_roa_prefix.roa_local_id and 
	rpki_roa.dir_id = rpki_dir.dir_id and
	rpki_roa.flags in (0x0144, 0x0104)	order by rpki_roa.asn `
	err = xormdb.XormEngine.Sql(sql).Find(&exports)
	if err != nil {
		fmt.Println("select roa failed, err is ", err)
		return
	}
	fmt.Println(len(exports))
	for i, _ := range exports {
		exports[i].Prefix = iputil.RtrFormatToIp(exports[i].AddressPrefix) + fmt.Sprintf("/%d", exports[i].PrefixLength)
		for key, value := range talMap {
			if strings.Contains(exports[i].DirName, key) {
				exports[i].Ta = value
				break
			}
		}
	}
	ext := osutil.Ext(file)
	switch ext {
	case ".json":
		json := jsonutil.MarshalJson(exports)
		data := []byte(json)
		err = ioutil.WriteFile(file, data, 0644)
	case ".csv":
		f, err := os.Create(file)
		if err != nil {
			fmt.Println("select roa failed, err is ", err)
			return
		}
		w := csv.NewWriter(f)
		for i, _ := range exports {
			w.Write([]string{strconv.Itoa(exports[i].Asn), exports[i].Prefix, strconv.Itoa(exports[i].MaxLength), exports[i].Ta})
		}
		w.Flush()
	default:
		fmt.Println("usage:\r\n1. get json: export ./file1.json\r\n2. get csv: export ./file1.csv")
		return
	}

	if err == nil {
		fmt.Println("write ", file, " ok")
	} else {
		fmt.Println("write ", file, " failed, err is ", err)
	}
}
