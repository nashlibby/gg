package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nashlibby/gk"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
)

type Data struct {
	AppName    string
	Port       string
	ModuleType string
	ModuleName string
	ModuleText string
	NeedAuth   bool
	ModelField ModelField
}

type ModelField struct {
	Name    string
	Model   string
	Declare string
	Type    string
	Comment string
	Json    string
	Output  bool
}

type Generator struct {
	ForceMode bool `json:"force_mode"`
	Data      Data `json:"data"`
}

func NewGenerator(data Data, forceMode bool) *Generator {
	if data.ModuleType != "project" {
		data.AppName = GetAppName()
	}
	return &Generator{
		ForceMode: forceMode,
		Data:      data,
	}
}

// 解析模板输出
func (g *Generator) ParseTemplate(tmplFile, outputFile string) {
	funcMap := template.FuncMap{
		"FirstUpper":  gk.FirstUpper,
		"FirstLetter": func(str string) string { return strings.ToLower(str[:1]) },
	}
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFS(Templates, "tmpl/"+tmplFile))
	var processed bytes.Buffer
	split := strings.Split(tmplFile, "/")
	fileName := split[len(split)-1]
	err := tmpl.ExecuteTemplate(&processed, fileName, g.Data)
	if err != nil {
		log.Fatalf("Error: Unable to parse data into template: %v\n", err)
	}
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Error: Could not format processed template: %v\n", err)
	}
	outputPath := "./" + outputFile
	if exist, _ := gk.FileExists(outputPath); !exist || g.ForceMode {
		log.Println("Generate file: ", outputPath)
		f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		defer f.Close()
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(outputPath), os.ModePerm)
			if err != nil {
				panic(err)
			}
			f, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer f.Close()
		}

		w := bufio.NewWriter(f)
		_, _ = w.WriteString(string(formatted))
		_ = w.Flush()
	}
}

// 获取应用名称
func GetAppName() string {
	file, err := os.OpenFile("./go.mod", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal("Error: Not find go.mod")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	firstLine, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal("Error: Not find app name")
	}
	r, _ := regexp.Compile(`module (.*)`)
	res := r.FindStringSubmatch(string(firstLine))
	if len(res) < 2 {
		log.Fatal("Error: Not find app name")
	}

	return strings.Trim(res[1], " ")
}

// 生成项目
func (g *Generator) GenProject() {
	log.Println("Run: go mod init " + g.Data.ModuleName)
	_ = exec.Command("bash", "-c", "go mod init "+g.Data.ModuleName).Run()
	g.ParseTemplate("main/main.tmpl", "main.go")
	g.GenApi("ping")
	g.GenApi("auth")
	g.GenCommon()
	g.GenMiddleware("cors")
	g.GenMiddleware("auth")
	g.GenMigrate()
	g.GenGitignore()
	g.GenConfig()
	g.GenDeploy()
	g.GenMakeFile()
}

// 生成模块
func (g *Generator) GenModule() {
	g.GenApi("api")
	g.GenModel()
	g.GenRepository()
	g.GenTransformer()
	g.GenService()
	g.GenRouter()
}

// 生成api
func (g *Generator) GenApi(tmplName string) {
	if tmplName == "ping" || tmplName == "auth" {
		g.ParseTemplate("api/"+tmplName+".tmpl", "app/api/"+tmplName+".go")
	} else {
		g.ParseTemplate("api/api.tmpl", "app/api/"+gk.Camel2Case(g.Data.ModuleName)+".go")
	}
}

// 生成model
func (g *Generator) GenModel() {
	outputFile := "app/internal/dao/model/" + gk.Camel2Case(g.Data.ModuleName) + ".go"
	g.ParseTemplate("dao/model.tmpl", outputFile)
	// 添加migrate
	migrateFilePath := "./app/migrate/migrate.go"
	content := "\t_ = common.DB.AutoMigrate(model." + gk.FirstUpper(g.Data.ModuleName) + "{})"
	if exists, _ := gk.StringExistsInFile(migrateFilePath, content); !exists {
		err := gk.InsertOneLineToFile(migrateFilePath, content, "func Run() {")
		if err != nil {
			log.Fatal(err)
		}
	}

	if exists, _ := gk.StringExistsInFile(migrateFilePath, "import"); !exists {
		includeText := fmt.Sprintf(`import (
	"%s/app/common"
	"%s/app/internal/dao/model"
)`, g.Data.AppName, g.Data.AppName)
		err := gk.InsertOneLineToFile(migrateFilePath, includeText, "package migrate")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if exists, _ := gk.StringExistsInFile(migrateFilePath, "app/common"); !exists {
			err := gk.InsertOneLineToFile(migrateFilePath, fmt.Sprintf(`"%s/app/common"`, g.Data.AppName), "import (")
			if err != nil {
				log.Fatal(err)
			}
		}
		if exists, _ := gk.StringExistsInFile(migrateFilePath, "app/internal/dao/model"); !exists {
			err := gk.InsertOneLineToFile(migrateFilePath, fmt.Sprintf(`"%s/app/internal/dao/model"`, g.Data.AppName), "import (")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// 生成repository
func (g *Generator) GenRepository() {
	g.ParseTemplate("dao/repository.tmpl", "app/internal/dao/repository/"+gk.Camel2Case(g.Data.ModuleName)+".go")
}

// 生成transformer
func (g *Generator) GenTransformer() {
	g.ParseTemplate("dao/transformer.tmpl", "app/internal/dao/transformer/"+gk.Camel2Case(g.Data.ModuleName)+".go")
}

// 生成service
func (g *Generator) GenService() {
	g.ParseTemplate("service/service.tmpl", "app/internal/service/"+gk.Camel2Case(g.Data.ModuleName)+".go")
}

// 生成middleware
func (g *Generator) GenMiddleware(tmplName string) {
	if tmplName == "blank" {
		g.ParseTemplate("middleware/blank.tmpl", "app/middleware/"+gk.Camel2Case(g.Data.ModuleName)+".go")
	} else {
		g.ParseTemplate("middleware/"+tmplName+".tmpl", "app/middleware/"+tmplName+".go")
	}
}

// 生成router
func (g *Generator) GenRouter() {
	g.ParseTemplate("router/router.tmpl", "app/router/"+gk.Camel2Case(g.Data.ModuleName)+".go")
	// 添加router
	content := "\trouter." + gk.FirstUpper(g.Data.ModuleName) + "Router(v1)"
	if exists, _ := gk.StringExistsInFile("./main.go", content); !exists {
		err := gk.InsertOneLineToFile("./main.go", content, "// 添加路由")
		if err != nil {
			log.Fatal(err)
		}
	}
	if exists, _ := gk.StringExistsInFile("./main.go", "app/router"); !exists {
		err := gk.InsertOneLineToFile("./main.go", fmt.Sprintf(`	"%s/app/router"`, g.Data.AppName), "import (")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 生成logic
func (g *Generator) GenLogic() {
	g.ParseTemplate("logic/logic.tmpl", "app/internal/logic/"+gk.Camel2Case(g.Data.ModuleName)+".go")
}

// 生成biz
func (g *Generator) GenBiz() {
	fileName := gk.Camel2Case(g.Data.ModuleName)
	g.ParseTemplate("biz/biz.tmpl", "app/internal/biz/"+fileName+"/"+fileName+".go")
}

// 生成common
func (g *Generator) GenCommon() {
	g.ParseTemplate("common/global.tmpl", "app/common/global.go")
	g.ParseTemplate("common/mysql.tmpl", "app/common/mysql.go")
	g.ParseTemplate("common/redis.tmpl", "app/common/redis.go")
}

// 生成migrate
func (g *Generator) GenMigrate() {
	g.ParseTemplate("migrate/migrate.tmpl", "app/migrate/migrate.go")
	g.ParseTemplate("migrate/seed.tmpl", "app/migrate/seed.go")
}

// 生成model field
func (g *Generator) GenField() {
	field := g.Data.ModelField
	modelFilePath := "app/internal/dao/model/" + gk.Camel2Case(g.Data.ModelField.Model) + ".go"
	if exists, err := gk.FileExists(modelFilePath); exists && err == nil {
		content := fmt.Sprintf("\t%s %s `gorm:\"type:%scomment:%s\" json:\"%s\"`", field.Name, field.Declare, field.Type, field.Comment, field.Json)
		err := gk.InsertOneLineToFile(modelFilePath, content, "gk.FieldTime", "before")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Error: not find model")
	}
	if field.Output {
		transformerFilePath := "app/internal/dao/transformer/" + gk.Camel2Case(g.Data.ModelField.Model) + ".go"
		if exists, err := gk.FileExists(transformerFilePath); exists && err == nil {
			content := fmt.Sprintf("\t%s %s `json:\"%s\"`", field.Name, field.Declare, field.Json)
			err := gk.InsertOneLineToFile(transformerFilePath, content, "CreatedAt", "before")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Error: not find transformer")
		}
	}
}

// 生成gitignore
func (g *Generator) GenGitignore() {
	if exist, _ := gk.FileExists(".gitignore"); !exist || g.ForceMode {
		_ = ioutil.WriteFile(".gitignore", []byte(fmt.Sprintf(`# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/

%s
logs
.idea
config.yaml
deploy.sh
`, g.Data.ModuleName)), os.ModePerm)
		log.Println("Generate file: ", "./.gitignore")
	}
}

// 生成config
func (g *Generator) GenConfig() {
	if exist, _ := gk.FileExists("config.yaml"); !exist || g.ForceMode {
		_ = ioutil.WriteFile("config.yaml", []byte(fmt.Sprintf(`app:
  name: %s
  port: %s
  mode: debug
  url: 127.0.0.1
db:
  host: 127.0.0.1
  port: 3306
  username: root
  password: root
  database: dbname
redis:
  host: 127.0.0.1
  port: 6379
  password:
  db: 1
`, g.Data.ModuleName, g.Data.Port)), os.ModePerm)
		log.Println("Generate file: ", "./config.yaml")
	}

}

// 生成deploy
func (g *Generator) GenDeploy() {
	if exist, _ := gk.FileExists("deploy.sh.example"); !exist || g.ForceMode {
		_ = ioutil.WriteFile("deploy.sh.example", []byte(fmt.Sprintf(`#!/bin/sh
git pull origin master
make build docker
cd docker-compose-path
docker-compose rm -sf app-%s
docker-compose up -d app-%s
`, g.Data.ModuleName, g.Data.ModuleName)), os.ModePerm)
		log.Println("Generate file: ", "./deploy.sh.example")
	}
}

// 生成makefile
func (g *Generator) GenMakeFile() {
	if exist, _ := gk.FileExists("Makefile"); !exist || g.ForceMode {
		_ = ioutil.WriteFile("Makefile", []byte(fmt.Sprintf(`.PHONY: build
build:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.cn go build -o %s *.go

.PHONY: docker
docker:
	@docker build -t app-%s:latest .

.PHONY: deploy
deploy:
	@./deploy.sh
`, g.Data.ModuleName, g.Data.ModuleName)), os.ModePerm)
		log.Println("Generate file: ", "./Makefile")
	}
}
