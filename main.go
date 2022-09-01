package main

import (
	"embed"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

//go:embed tmpl
var Templates embed.FS

func main() {

	app := &cli.App{
		Name:                   "A cli tool for generating go code with gin and gorm",
		Usage:                  "Fast generate curd api code",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "project",
				Usage: "Init project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "port", Value: "8080", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The project name is required, please use the -a or --name flag to specify the project name")
					}
					data := Data{
						ModuleType: "project",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
						Port:       cCtx.String("port"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenProject()
					return nil
				},
			},
			{
				Name:  "module",
				Usage: "Generate module code",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The module name is required, please use the -a or --name flag to specify the module name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenModule()
					return nil
				},
			},
			{
				Name:  "api",
				Usage: "Generate api code",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The api name is required, please use the -a or --name flag to specify the api name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenApi("api")
					return nil
				},
			},
			{
				Name:  "model",
				Usage: "Generate model code",
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
				Name:  "repository",
				Usage: "Generate repository code",
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
				Name:  "transformer",
				Usage: "Generate transformer code",
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
				Name:  "logic",
				Usage: "Generate logic code",
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
				Name:  "service",
				Usage: "Generate service code",
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
				Name:  "middleware",
				Usage: "Generate middleware code",
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
				Name:  "router",
				Usage: "Generate router code",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"a"}},
					&cli.StringFlag{Name: "text", Aliases: []string{"t"}, DefaultText: ""},
					&cli.BoolFlag{Name: "force", Value: false, Aliases: []string{"f"}},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") == "" {
						log.Fatal("Error: The router name is required, please use the -a or --name flag to specify the router name")
					}
					data := Data{
						ModuleType: "module",
						ModuleName: cCtx.String("name"),
						ModuleText: cCtx.String("text"),
					}
					NewGenerator(data, cCtx.Bool("force")).GenRouter()
					return nil
				},
			},
			{
				Name:  "biz",
				Usage: "Generate biz code",
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	} else {
		log.Println("DONE!")
	}
}
