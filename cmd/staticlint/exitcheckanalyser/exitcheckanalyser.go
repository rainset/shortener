// Package exitcheckanalyser является кастомным анализатором os.Exit в функции main

package exitcheckanalyser

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

// хранение позиции
type Position struct {
	Name string
	Pos  token.Pos
	End  token.Pos
}

// запуск
func run(pass *analysis.Pass) (interface{}, error) {
	//fset := token.NewFileSet()
	//pass.Reportf(1, "os.Exit declaration")
	for _, file := range pass.Files {

		posPackage := make(map[string]Position)
		posFunc := make(map[string]Position)
		//posCallExpr := make(map[string]Position)

		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(fileNode ast.Node) bool {

			switch x := fileNode.(type) {
			case *ast.File:
				///ast.Print(fset, fileNode.Pos())
				//log.Println("package:", x.Name.Name, x.Pos(), x.End())

				posPackage[x.Name.Name] = Position{Name: x.Name.Name, Pos: x.Pos(), End: x.End()}

			case *ast.CallExpr:
				//log.Println("CallExpr", x.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name, x.Fun.(*ast.SelectorExpr).Sel.Name)
				name := x.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name
				selName := x.Fun.(*ast.SelectorExpr).Sel.Name
				//posCallExpr[name] = Position{Name: name, Pos: x.Pos(), End: x.End()}

				if name != "os" || selName != "Exit" {
					return true
				}

				if _, ok := posPackage["main"]; !ok {
					return true
				}

				if _, ok := posFunc["main"]; !ok {
					return true
				}

				if posPackage["main"].Pos < posFunc["main"].Pos &&
					posPackage["main"].End >= posFunc["main"].End &&
					posFunc["main"].Pos < x.Pos() &&
					posFunc["main"].End >= x.End() {

					//log.Println("posPackage", posPackage)
					//log.Println("FuncDecl", posFunc)
					pass.Reportf(x.Pos(), "os.Exit declaration")
				}

			case *ast.FuncDecl:
				//log.Println("FuncDecl", x.Name.Name, x.Pos(), x.End())
				posFunc[x.Name.Name] = Position{Name: x.Name.Name, Pos: x.Pos(), End: x.End()}
				//ast.Print(fset, x)
			}

			return true
		})
	}
	return nil, nil
}
