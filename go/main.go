package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	tmpl = "template/"
	md   = "markdown/"
	out  = "public/html/"
	act  = "activity/"
	des  [][]string
	err  error
)

func main() {

	//init

	os.RemoveAll(out + "activity")
	os.Remove(out + "index.html")
	os.Remove(out + "list.html")
	os.Remove("./public/sitemap.xml")

	// checknil(err)
	err = os.Mkdir(out+act, os.ModePerm)
	checknil(err)

	//activity
	mds, err := ioutil.ReadDir(md + "活動/")
	checknil(err)

	for i, el := range mds {
		tmp := out + act + el.Name() + "/"
		mdfile := md + "活動/" + el.Name() + "/index.md"
		value := getvalue(mdfile)
		file, err := template.New("tmpl").
			Funcs(template.FuncMap{"kaiyo": func() string { return "" }}).
			ParseFiles(tmpl+"activity.html", mdfile)
		checknil(err)
		err = os.Mkdir(tmp, os.ModePerm)
		checknil(err)
		put, err := os.Create(tmp + "/index.html")
		checknil(err)
		value["ISO"] = strings.Replace(el.Name(), " ", "-", 2)
		err = file.Execute(put, value)
		checknil(err)
		des = append([][]string{{el.Name(), value["標題"], value["簡介"]}}, des...)
		copy, err := ioutil.ReadDir(md + "活動/" + el.Name() + "/")
		checknil(err)
		for _, el0 := range copy {
			if el0.Name() != "index.md" {
				src, err := os.Open(md + "活動/" + el.Name() + "/" + el0.Name())
				checknil(err)
				defer src.Close()
				dis, err := os.Create(tmp + "/" + el0.Name())
				checknil(err)
				defer dis.Close()
				_, err = io.Copy(dis, src)
				checknil(err)
			}
		}
		fmt.Printf("activity done! %d / %d \n", i+1, len(mds))
	}

	//home
	file, err := template.ParseFiles(tmpl + "index.html")
	checknil(err)
	value := getvalue(md + "首頁.gtpl")
	put, err := os.Create(out + "index.html")
	checknil(err)
	exe := struct {
		Home map[string]string
		News [][]string
	}{
		Home: value,
		News: des,
	}
	err = file.Execute(put, exe)
	checknil(err)
	fmt.Printf("index done! \n")

	//list
	file, err = template.ParseFiles(tmpl + "list.html")
	checknil(err)
	put, err = os.Create(out + "list.html")
	checknil(err)
	err = file.Execute(put, des)
	checknil(err)
	fmt.Printf("list done! \n")

	//sitemap
	file, err = template.ParseFiles(tmpl + "sitemap.xml")
	checknil(err)
	put, err = os.Create("public/sitemap.xml")
	checknil(err)
	err = file.Execute(put, des)
	checknil(err)
	fmt.Printf("sitemap done! \n")

}
