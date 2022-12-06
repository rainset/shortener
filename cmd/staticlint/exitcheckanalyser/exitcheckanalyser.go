// Package exitcheckanalyser является кастомным анализатором os.Exit в функции main

package exitcheckanalyser

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	//fset := token.NewFileSet()
	//pass.Reportf(1, "os.Exit declaration")
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(fileNode ast.Node) bool {

			switch x := fileNode.(type) {
			case *ast.File:

				//log.Println("package:", x.Name.Name)
				for i := 0; i < len(x.Decls); i++ {
					if funcDecl, ok := x.Decls[i].(*ast.FuncDecl); ok {
						//log.Println("funcDecl:", funcDecl.Body)
						for i := 0; i < len(funcDecl.Body.List); i++ {
							if exprStmt, ok := funcDecl.Body.List[i].(*ast.ExprStmt); ok {
								//log.Println("exprStmt:", exprStmt.X.(*ast.CallExpr).Fun)
								if callExpr, ok := exprStmt.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr); ok {
									//log.Println("callExpr11:", x.Name.Name, callExpr.X.(*ast.Ident).Name, callExpr.Sel.Name)

									//log.Println("callExpr", exprStmt.X.(*ast.CallExpr))

									if callExpr.X.(*ast.Ident).Name == "os" &&
										callExpr.Sel.Name == "Exit" &&
										x.Name.Name == "main" {
										//log.Println(x.Pos(), callExpr.X.(*ast.Ident).Name, callExpr.Sel.Name, x.Name.Name)

										//log.Println("os.Exit declaration")
										pass.Reportf(callExpr.Pos(), "os.Exit declaration")
									}

								}
							}
						}
					}
				}
			case *ast.CallExpr:
			case *ast.FuncDecl:
			}

			return true
		})
	}
	return nil, nil
}
