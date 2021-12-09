package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// proto的message的结构信息
type ProtoMessageStructInfo struct {
	// *.proto的文件名(不包含.proto)
	protoName   string
	// message结构名
	messageName string
	// 注释关键字
	keyComment  string
	// 注释关键字可以对应一个值
	keyCommentValue string
	// 排除keyComment后的注释
	normalComment string
	// *pb.go的package name
	pbPackageName string
	// message的结构信息
	structType  *ast.StructType
}

// 代码模板
type CodeTemplate struct {
	// 注释关键字,如@Player
	KeyComment string

	// 生成文件名
	OutFile string

	/*
	package game
	import "github.com/fish-tennis/gserver/pb"
	*/
	// 文件头,用[]string,为了解决code_templates.json里不方便写换行的问题
	Header []string

	/*
	@Player对应的函数模板:
	func (this *Player) Send{MessageName}(packet *pb.{MessageName}) bool {
		return this.Send(Cmd(pb.Cmd{ProtoFileName}_Cmd_{MessageName}), packet)
	}
	@Server对应的函数模板
	func SendPacket{MessageName}(conn Connection, packet *pb.{MessageName}) bool {
		return conn.Send(Cmd(pb.Cmd{ProtoFileName}_Cmd_{MessageName}), packet)
	}
	*/
	// 函数替换模板
	FuncTemplate []string

	// 文件尾
	Tail []string
}

type CodeTemplates []*CodeTemplate

type ParserResult struct {
	codeTemplates []*CodeTemplate
	protoMap map[string][]*ProtoMessageStructInfo // key:protoName
}

func (this *ParserResult) GetCodeTemplate(key string) *CodeTemplate {
	for _,v := range this.codeTemplates {
		if v.KeyComment == key {
			return v
		}
	}
	return nil
}

// 解析protoc-gen-go生成的*pb.go代码
// 参考github.com/favadi/protoc-go-inject-tag
// protoc-go-inject-tag只能对Message的字段(field)进行处理,不能完全满足我们的需求,我们希望直接对Message(struct)的注释进行解析
// 而且golang的反射只能获取field的tag,struct自身没有tag,因此如果我们希望对struct进行特殊标记,生成的*pb.go代码,也无法处理struct
// 所以我们的解决方案是,直接利用golang的parser接口,解析出struct的注释信息,根据在注释里插入的关键字,生成辅助代码,应用层调用生成的辅助代码时,
// 由于不需要再进行反射操作,性能也没有损失
func ParseProtoCode(protoCodeFile string, parserResult *ParserResult) {
	//log.Printf("protoCodeFile:%v", protoCodeFile)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, protoCodeFile, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, decl := range f.Decls {
		// check if is generic declaration
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var typeSpec *ast.TypeSpec
		for _, spec := range genDecl.Specs {
			if ts, tsOK := spec.(*ast.TypeSpec); tsOK {
				typeSpec = ts
				break
			}
		}

		// skip if can't get type spec
		if typeSpec == nil {
			continue
		}
		//println(fmt.Sprintf("typeSpec:%v", typeSpec))

		// not a struct, skip
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		//println(fmt.Sprintf("struct doc:%v", genDecl.Doc))
		// struct的注释在genDecl.Doc里
		if genDecl.Doc != nil {
			keyChecker := make(map[string]struct{})
			for i,structComment := range genDecl.Doc.List {
				comment := strings.TrimPrefix(structComment.Text, "//")
				comment = strings.TrimSpace(comment)
				commentValue := ""
				normalComment := ""
				var codeTemplate *CodeTemplate
				for _,template := range parserResult.codeTemplates {
					lowerComment := strings.ToLower(comment)
					lowerKey := strings.ToLower(template.KeyComment)
					// 不区分大小写
					if lowerComment == lowerKey {
						codeTemplate = template
					} else if strings.HasPrefix(lowerComment, lowerKey) && strings.Contains(lowerComment,":") {
						kv := strings.Split(comment,":")
						if len(kv) == 2 && kv[1] != "" {
							codeTemplate = template
							commentValue = strings.TrimSpace(kv[1])
						}
					}
					if codeTemplate != nil && normalComment == "" {
						for j,v := range genDecl.Doc.List {
							if j != i {
								if normalComment != "" {
									normalComment += "\n"
								}
								normalComment += v.Text
							}
						}
					}
				}
				if codeTemplate != nil {
					// 排重
					if _,ok := keyChecker[comment]; ok {
						continue
					}
					structInfo := &ProtoMessageStructInfo{
						protoName:   path.Base(path.Clean(strings.Replace(protoCodeFile,"\\","/",-1))),
						messageName: typeSpec.Name.Name,
						keyComment:  codeTemplate.KeyComment,
						keyCommentValue: commentValue,
						normalComment: normalComment,
						pbPackageName: f.Name.Name,
						structType:  structDecl,
					}
					structInfoList := parserResult.protoMap[structInfo.protoName]
					if structInfoList == nil {
						structInfoList = make([]*ProtoMessageStructInfo, 0)
					}
					structInfoList = append(structInfoList, structInfo)
					parserResult.protoMap[structInfo.protoName] = structInfoList
					keyChecker[comment] = struct{}{}
					println(fmt.Sprintf("%v KeyComment:%v value:%v", structInfo.messageName, structInfo.keyComment, structInfo.keyCommentValue))
				}
			}
		}
	}
}

// 解析*pb.go文件
func ParseFiles(pbGoFilePattern string, codeTemplatesConfig string) {
	codeTemplates := initCodeTemplatesConfig(codeTemplatesConfig)
	parserResult := &ParserResult{
		codeTemplates: codeTemplates,
		protoMap: map[string][]*ProtoMessageStructInfo{},
	}
	files, err := filepath.Glob(pbGoFilePattern)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file count:%v", len(files))
	for _, path := range files {
		finfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		if finfo.IsDir() {
			continue
		}

		// It should end with ".pb.go" at a minimum.
		if !strings.HasSuffix(strings.ToLower(finfo.Name()), ".pb.go") {
			continue
		}

		ParseProtoCode(path, parserResult)
	}
	for _,codeTemplate := range parserResult.codeTemplates {
		generateCode(parserResult, codeTemplate.KeyComment)
	}
}

// 生成相应的辅助代码
func generateCode(parserResult *ParserResult, key string) {
	codeTemplate := parserResult.GetCodeTemplate(key)
	builder := strings.Builder{}
	builder.WriteString(strings.Join(codeTemplate.Header, "\n"))
	for _,structInfoList := range parserResult.protoMap {
		for _,structInfo := range structInfoList {
			if structInfo.keyComment != codeTemplate.KeyComment {
				continue
			}
			protoFileName := structInfo.protoName[:strings.Index(structInfo.protoName,".pb.go")]
			protoName := protoFileName
			// 首字母大写
			ProtoName := strings.ToUpper(protoFileName[:1]) + protoFileName[1:]
			messageName := structInfo.messageName
			funcStr := strings.Join(codeTemplate.FuncTemplate, "\n")
			// 替换掉代码模板中的关键字
			funcStr = strings.ReplaceAll(funcStr, "{MessageName}", messageName)
			funcStr = strings.ReplaceAll(funcStr, "{protoName}", protoName)
			funcStr = strings.ReplaceAll(funcStr, "{ProtoName}", ProtoName)
			funcStr = strings.ReplaceAll(funcStr, "{PackageName}", structInfo.pbPackageName)
			funcStr = strings.ReplaceAll(funcStr, "{Value}", structInfo.keyCommentValue)
			funcStr = strings.ReplaceAll(funcStr, "{Comment}", structInfo.normalComment)
			builder.WriteString(funcStr)
			builder.WriteString("\n")
		}
	}
	builder.WriteString(strings.Join(codeTemplate.Tail, "\n"))
	os.WriteFile(codeTemplate.OutFile, ([]byte)(builder.String()), 0644)
}

func initCodeTemplatesConfig(config string) []*CodeTemplate {
	fileData,err := os.ReadFile(config)
	if err != nil {
		panic("read config file err")
	}
	var codeTemplates CodeTemplates
	err = json.Unmarshal(fileData, &codeTemplates)
	if err != nil {
		panic(err)
	}
	return codeTemplates
}

func main() {
	var codeTemplatesConfigFile,inputFiles string
	flag.StringVar(&inputFiles, "input", "", "pattern to match input file(s)")
	flag.StringVar(&codeTemplatesConfigFile, "config", "", "code template config")
	flag.Parse()

	if len(inputFiles) == 0 {
		log.Fatal("no input file")
	}
	ParseFiles(inputFiles, codeTemplatesConfigFile)
}