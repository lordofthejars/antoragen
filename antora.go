package main

import (
	"bufio"
	html "html/template"
	"os"
	"path"
	"text/template"

	"github.com/gobuffalo/packr"
)

var box = packr.NewBox("./tmpl")
var sitebox = packr.NewBox("./tmpl-site")

func generateSite(outputDir, projectName, repo, public, startRepo string) error {
	err := os.Chdir(outputDir)

	if err != nil {
		return err
	}

	err = renderTemplate("site.yml.tmpl", path.Join(outputDir, "site.yml"), map[string]string{"projectName": projectName, "repo": repo, "public": public, "startRepo": startRepo}, sitebox)

	if err != nil {
		return err
	}

	docsFolder := path.Join(outputDir, "docs")
	err = os.MkdirAll(docsFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = writeFile("nojekyll", path.Join(docsFolder, ".nojekyll"), sitebox)

	if err != nil {
		return err
	}

	supplementalUiFolder := path.Join(outputDir, "supplemental-ui")

	err = os.MkdirAll(supplementalUiFolder, os.ModePerm)

	if err != nil {
		return err
	}

	imgFolder := path.Join(supplementalUiFolder, "img")

	err = os.MkdirAll(imgFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = writeBinaryFile("favicon.ico", path.Join(imgFolder, "favicon.ico"), sitebox)

	if err != nil {
		return err
	}

	partialsFolder := path.Join(supplementalUiFolder, "partials")

	err = os.MkdirAll(partialsFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = writeFile("head-meta.hbs", path.Join(partialsFolder, "head-meta.hbs"), sitebox)

	if err != nil {
		return err
	}

	err = renderTemplate("header-content.hbs.tmpl", path.Join(partialsFolder, "header-content.hbs"), map[string]string{"projectName": projectName, "repo": repo}, sitebox)

	if err != nil {
		return err
	}

	return nil

}

func generateDoc(outputDir, projectName string) error {

	err := os.Chdir(outputDir)

	if err != nil {
		return err
	}

	docsFolder := path.Join(outputDir, "docs")

	err = os.MkdirAll(docsFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = renderTemplate("antora.yml.tmpl", path.Join(docsFolder, "antora.yml"), map[string]string{"projectName": projectName}, box)

	if err != nil {
		return err
	}

	modulesFolder := path.Join(docsFolder, "modules")

	err = os.MkdirAll(modulesFolder, os.ModePerm)

	if err != nil {
		return err
	}

	rootFolder := path.Join(modulesFolder, "ROOT")

	err = os.MkdirAll(rootFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = writeFile("globalattributes.adoc", path.Join(rootFolder, "_attributes.adoc"), box)

	if err != nil {
		return err
	}

	err = renderTemplate("nav.adoc.tmpl", path.Join(rootFolder, "nav.adoc"), map[string]string{"projectName": projectName}, box)

	if err != nil {
		return err
	}

	assetsFolder := path.Join(rootFolder, "assets")

	err = os.MkdirAll(assetsFolder, os.ModePerm)

	if err != nil {
		return err
	}

	imagesFolder := path.Join(assetsFolder, "images")

	err = os.MkdirAll(imagesFolder, os.ModePerm)

	if err != nil {
		return err
	}

	pagesFolder := path.Join(rootFolder, "pages")

	err = os.MkdirAll(pagesFolder, os.ModePerm)

	if err != nil {
		return err
	}

	err = renderTemplate("index.adoc.tmpl", path.Join(pagesFolder, "index.adoc"), map[string]string{"projectName": projectName}, box)

	if err != nil {
		return err
	}

	err = writeFile("moduleattributes.adoc", path.Join(pagesFolder, "_attributes.adoc"), box)

	if err != nil {
		return err
	}

	err = writeFile("section.adoc", path.Join(pagesFolder, "section.adoc"), box)

	if err != nil {
		return err
	}

	err = writeFile("anothersection.adoc", path.Join(pagesFolder, "anothersection.adoc"), box)

	if err != nil {
		return err
	}

	return nil
}

func writeBinaryFile(file, outputFile string, box packr.Box) error {

	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()

	contentFile, err := box.MustBytes(file)

	if err != nil {
		return err
	}

	_, err = f.Write(contentFile)

	if err != nil {
		return err
	}

	f.Sync()

	return nil

}

func writeFile(file, outputFile string, box packr.Box) error {
	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()

	contentFile, err := box.MustString(file)

	if err != nil {
		return err
	}

	_, err = f.WriteString(contentFile)

	if err != nil {
		return err
	}

	f.Sync()

	return nil

}

func renderHtmlTemplate(tmplName, outputFile string, p map[string]string, box packr.Box) error {
	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()

	templateFile, err := box.MustString(tmplName)

	if err != nil {
		return err
	}

	templates, _ := html.New(tmplName).Parse(templateFile)
	w := bufio.NewWriter(f)

	templates.Execute(w, p)
	w.Flush()

	return nil
}

func renderTemplate(tmplName, outputFile string, p map[string]string, box packr.Box) error {
	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()

	templateFile, err := box.MustString(tmplName)

	if err != nil {
		return err
	}

	templates, _ := template.New(tmplName).Parse(templateFile)
	w := bufio.NewWriter(f)

	templates.Execute(w, p)
	w.Flush()

	return nil
}
