package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"baliance.com/gooxml/document"
	"baliance.com/gooxml/schema/soo/wml"
)

type TreeStorage struct {
	Value map[string]any
	mutex sync.Mutex
}

func (t *TreeStorage) addElem(key string, value any) {
	t.mutex.Lock()
	t.Value[key] = value
	t.mutex.Unlock()
}


func main() {
	var wg sync.WaitGroup
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current directory:", dir)
	storage := &TreeStorage{Value: map[string]any{}, mutex: sync.Mutex{}}

	files, _ := os.ReadDir(dir)
	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		if strings.HasPrefix(file.Name(), ".") || file.Name() == "structure.docx" || strings.HasSuffix(file.Name(), ".exe") {
			continue
		}
		if file.IsDir() {
			wg.Add(1)
			go recursiveFileParser(fpath, &wg, storage)
		} else {
			b, err := os.ReadFile(fpath)
			if err == nil {
				file_text := string(b)
				storage.addElem(file.Name(), file_text)
			}
		}
	}
	wg.Wait()
	printTreeStructure("", storage.Value)
	saveStructureToDocx(&storage.Value)
}

func recursiveFileParser(dir string, wg *sync.WaitGroup, storage *TreeStorage) {
	files, _ := os.ReadDir(dir)
	res := map[string]any{}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasSuffix(file.Name(), ".exe") {
			continue
		}
		fpath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			res[file.Name()] = recursiveSyncFileParser(fpath)
			continue
		}
		b, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		file_text := string(b)
		res[file.Name()] = file_text
	}
	storage.addElem(filepath.Base(dir), res)
	wg.Done()
}

func recursiveSyncFileParser(dir string) map[string]any {
	files, _ := os.ReadDir(dir)
	res := map[string]any{}
	for _, file := range files {
		fpath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			res[file.Name()] = recursiveSyncFileParser(fpath)
			continue
		}
		b, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}
		file_text := string(b)
		res[file.Name()] = file_text
	}
	return res
}

func printTreeStructure(prefix string, tree map[string]any) {
	p := prefix + "|---"
	
	for k, v := range tree {
		fmt.Printf("%s%s\n", p, k)
		value, ok := v.(map[string]any)
		if ok {
			printTreeStructure(p, value)
		}
	}
}

func saveStructureToDocx(tree *map[string]any) {
	doc := document.New()
	para := doc.AddParagraph()
	run := para.AddRun()
	para.SetStyle("Title")
	run.AddText("Files Tree Structure")
	addTreeStructure("", tree, doc)

	doc.AddParagraph().Properties().AddSection(wml.ST_SectionMarkNextPage)

	para = doc.AddParagraph()
	para.SetStyle("Title")
	para.AddRun().AddText("Files Content")
	addFilesContent("", tree, doc)
	doc.SaveToFile("structure.docx")
}

func addTreeStructure(prefix string, tree *map[string]any, d *document.Document) {
	p := prefix + "|---"
	
	for k, v := range *tree {
		s := fmt.Sprintf("%s%s\n", p, k)
		para := d.AddParagraph()
		run := para.AddRun()
		run.AddText(s)
		value, ok := v.(map[string]any)
		if ok {
			addTreeStructure(p, &value, d)
		}
	}
}

func addFilesContent(name_prefix string, tree *map[string]any, d *document.Document) {
	for k, v := range *tree {
		var title string
		if name_prefix != "" {
			title = fmt.Sprintf("%s/%s\n", name_prefix, k)
		} else {
			title = fmt.Sprintf("%s\n", k)
		}
		para := d.AddParagraph()
		para.Properties().SetHeadingLevel(2)
		para.AddRun().AddText(title)
		value, ok := v.(map[string]any)
		if ok {
			addFilesContent(title, &value, d)
			continue
		}
		lines := strings.Split(v.(string), "\n")
		for _, line := range lines {
			para := d.AddParagraph()
			run := para.AddRun()
			run.AddText(line)
			run.Properties().SetFontFamily("Courier New")
			run.Properties().SetSize(9)
		}
	}
}