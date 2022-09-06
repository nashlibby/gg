package main

import (
	"embed"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
)

//go:embed tmpl
var Templates embed.FS

func main() {

	app := &cli.App{
		Name:                   "A cli tool for generating go code with gin and gorm",
		Usage:                  "Fast generate curd api code",
		UsageText:              "gg command [command options] [arguments...]",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:      "project",
				Usage:     "Init project",
				UsageText: "eg. gg project -a cms -t 内容管理系统 -p 8088 -s -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "port", Value: "8080", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
					&cli.BoolFlag{Name: "swagger", Value: false, Aliases: []string{"s"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The project name is required, please use the -a or --name flag to specify the project name")
					}
					data := Data{
						ModuleType: "project",
						AppName:    cCtx.String("name"),
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
						Port:       cCtx.String("port"),
						UseSwagger: cCtx.Bool("swagger"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenProject()
					log.Println("Run: go mod tidy")
					_ = exec.Command("bash", "-c", "go mod tidy").Run()
					if cCtx.Bool("swagger") {
						log.Println("Run: swag init")
						_ = exec.Command("bash", "-c", "swag init").Run()
					}
					return nil
				},
			},
			{
				Name:      "module",
				Usage:     "Generate module code",
				UsageText: "eg. gg module -a user -t 用户 -u -s -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
					&cli.BoolFlag{Name: "auth", Value: false, Aliases: []string{"u"}},
					&cli.BoolFlag{Name: "swagger", Value: false, Aliases: []string{"s"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The module name is required, please use the -a or --name flag to specify the module name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
						NeedAuth:   cCtx.Bool("auth"),
						UseSwagger: cCtx.Bool("swagger"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenModule()
					log.Println("Run: go mod tidy")
					_ = exec.Command("bash", "-c", "go mod tidy").Run()
					if cCtx.Bool("swagger") {
						log.Println("Run: swag init")
						_ = exec.Command("bash", "-c", "swag init").Run()
					}
					return nil
				},
			},
			{
				Name:      "api",
				Usage:     "Generate api code",
				UsageText: "gg api -a user -t 用户 -s -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
					&cli.BoolFlag{Name: "swagger", Value: false, Aliases: []string{"s"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The api name is required, please use the -a or --name flag to specify the api name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
						UseSwagger: cCtx.Bool("swagger"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenApi("api")
					if cCtx.Bool("swagger") {
						log.Println("Run: swag init")
						_ = exec.Command("bash", "-c", "swag init").Run()
					}
					return nil
				},
			},
			{
				Name:      "model",
				Usage:     "Generate model code",
				UsageText: "gg model -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The model name is required, please use the -a or --name flag to specify the model name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenModel()
					return nil
				},
			},
			{
				Name:      "repository",
				Usage:     "Generate repository code",
				UsageText: "gg repository -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The repository name is required, please use the -a or --name flag to specify the repository name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenRepository()
					return nil
				},
			},
			{
				Name:      "transformer",
				Usage:     "Generate transformer code",
				UsageText: "gg transformer -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The transformer name is required, please use the -a or --name flag to specify the transformer name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenTransformer()
					return nil
				},
			},
			{
				Name:      "logic",
				Usage:     "Generate logic code",
				UsageText: "gg logic -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The logic name is required, please use the -a or --name flag to specify the logic name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenLogic()
					return nil
				},
			},
			{
				Name:      "service",
				Usage:     "Generate service code",
				UsageText: "gg logic -a service -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The service name is required, please use the -a or --name flag to specify the service name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenService()
					return nil
				},
			},
			{
				Name:      "middleware",
				Usage:     "Generate middleware code",
				UsageText: "gg middleware -a auth -t 认证 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The middleware name is required, please use the -a or --name flag to specify the middleware name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenMiddleware("blank")
					return nil
				},
			},
			{
				Name:      "router",
				Usage:     "Generate router code",
				UsageText: "gg router -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
					&cli.BoolFlag{Name: "auth", Value: false, Aliases: []string{"u"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The router name is required, please use the -a or --name flag to specify the router name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
						NeedAuth:   cCtx.Bool("auth"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenRouter()
					return nil
				},
			},
			{
				Name:      "biz",
				Usage:     "Generate biz code",
				UsageText: "gg biz -a user -t 用户 -f",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The biz name is required, please use the -a or --name flag to specify the biz name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenBiz()
					return nil
				},
			},
			{
				Name:  "field",
				Usage: "Generate field of model",
				UsageText: `eg. gg field -m User -a UserName -d string -t "varchar(30);not null;" -c "用户名" -j user_name -o
eg. gg field -m User -a Age -d uint8 -t "tinyint(3);unsigned;not null;default:0;" -c "年龄" -j age`,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "model", Aliases: []string{"m"}},
					&cli.StringFlag{Name: "declare", Aliases: []string{"d"}},
					&cli.StringFlag{Name: "type", Aliases: []string{"t"}},
					&cli.StringFlag{Name: "comment", Aliases: []string{"c"}},
					&cli.StringFlag{Name: "json", Aliases: []string{"j"}},
					&cli.BoolFlag{Name: "output", Value: false, Aliases: []string{"o"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The name of model field is required, please use the -a or --name flag to specify the name of model field")
					}

					if cCtx.String("model") == "" {
						log.Fatal("Error: The name of model is required, please use the -m or --model flag to specify the name of model")
					}

					if cCtx.String("declare") == "" {
						log.Fatal("Error: The field declare is required, please use the -d or --declare flag to specify the model field declare")
					}

					if cCtx.String("type") == "" {
						log.Fatal("Error: The type of field is required, please use the -t or --type flag to specify the type of field")
					}

					if cCtx.String("comment") == "" {
						log.Fatal("Error: The comment of field is required, please use the -c or --comment flag to specify the comment of field")
					}

					if cCtx.String("json") == "" {
						log.Fatal("Error: The json declare of field is required, please use the -j or --json flag to specify the json declare of field")
					}

					data := Data{
						ModuleType: "field",
						ModelField: ModelField{
							Name:    cCtx.String("name"),
							Model:   cCtx.String("model"),
							Declare: cCtx.String("declare"),
							Type:    cCtx.String("type"),
							Comment: cCtx.String("comment"),
							Json:    cCtx.String("json"),
							Output:  cCtx.Bool("output"),
						},
					}

					NewGenerator(data, false).GenField()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	} else {
		log.Println("DONE!")
	}
}
