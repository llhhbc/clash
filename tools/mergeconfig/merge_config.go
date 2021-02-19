package main

import (
	"flag"
	"fmt"
	"github.com/Dreamacro/clash/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

var (
	srcFile = flag.String("src", "", "src config file names. split by ,")
	desFile = flag.String("des", "", "empty for stdout. ")
)

func main()  {
	flag.Parse()

	srcFiles := strings.Split(*srcFile, ",")

	if *srcFile == "" || len(srcFiles) == 0 {
		log.Fatal("get empty src files. ")
	}

	msg, err := ioutil.ReadFile(srcFiles[0])
	if err != nil {
		log.Fatalf("get file %s data failed %v. ", srcFiles[0], err)
	}
	resConfig, err := config.UnmarshalRawConfig(msg)
	if err != nil {
		log.Fatalf("unmarshal first file failed %v. ", err)
	}

	for i:=1;i<len(srcFiles);i++ {
		msg, err = ioutil.ReadFile(srcFiles[i])
		if err != nil {
			log.Fatalf("get file %s data failed %v. ", srcFiles[i], err)
		}
		cur, err := config.UnmarshalRawConfig(msg)
		if err != nil {
			log.Fatalf("unmarshal file %s failed %v. ", srcFiles[i], err)
		}
		for _, p := range cur.Proxy {
			resConfig.Proxy = append(resConfig.Proxy, p)
		}
		for _, pg := range cur.ProxyGroup {
			resConfig.ProxyGroup = append(resConfig.ProxyGroup, pg)
		}
		for _, r := range cur.Rule {
			resConfig.Rule = append(resConfig.Rule, r)
		}
	}

	m, _ := yaml.Marshal(resConfig)
	if *desFile == "" {
		fmt.Println(string(m))
	} else {
		err = ioutil.WriteFile(*desFile, m, 0755)
		if err != nil {
			log.Fatalf("write file failed %v. ", err)
		}
	}

}
