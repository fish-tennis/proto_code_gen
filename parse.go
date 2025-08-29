package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// proto的message的结构信息
type ProtoMessageStructInfo struct {
	// *.pb.go的文件名
	protoName string
	// message结构名
	MessageName string
	// 注释关键字
	keyComment string
	// 注释关键字可以对应一个值
	Value string
	// 排除keyComment后的注释
	Comment string
	// *pb.go的package name
	PackageName string
	// message的结构信息
	structType *ast.StructType

	Fields []*ReaderField
}

type ReaderField struct {
	FieldName string
	FieldType string
	// 如果是Pointer,移除前面的星号(如 *string -> string)
	FieldTypeWithoutPtr string
	// slice和map的elem类型名(移除了前面的星号)
	ElemTypeName   string
	MapKeyTypeName string
	IsStruct       bool // 子message,如 *Item
	IsNormalSlice  bool // 普通slice,如[]int32
	IsPtrSlice     bool // pointer slice, 如[]*Item
	IsNormalMap    bool // 普通map,如map[string]int32
	IsPtrMap       bool // pointer map,如map[int32]*Item
}

// 自定义tag代码模板
type CodeTemplate struct {
	// 注释关键字,如@Player
	KeyComment string `yaml:"KeyComment"`

	// 模板文件
	Template string `yaml:"Template"`

	// 生成文件名
	OutFile string `yaml:"OutFile"`
}

// Reader代码模板
type ReaderCodeTemplate struct {
	// 模板文件
	Template string `yaml:"Template"`

	// 生成目录
	OutDir string `yaml:"OutDir"`

	// 处理哪些文件,支持正则
	FileFilter []string `yaml:"FileFilter"`

	// 处理哪些message,支持正则
	MessageFilter []string `yaml:"MessageFilter"`

	// 是否使用proto2
	ProtoV2 bool `yaml:"ProtoV2"`
}

type CommandMapping struct {
	OutFile string `yaml:"OutFile"`
}

func (this *ReaderCodeTemplate) MatchFile(fileName string) bool {
	for _, filter := range this.FileFilter {
		if ok, _ := regexp.MatchString(filter, fileName); ok {
			return true
		}
	}
	return false
}

func (this *ReaderCodeTemplate) MatchMessage(messageName string) bool {
	for _, filter := range this.MessageFilter {
		if ok, _ := regexp.MatchString(filter, messageName); ok {
			return true
		}
	}
	return false
}

type CodeTemplates struct {
	Code           []*CodeTemplate     `yaml:"Code"`
	Reader         *ReaderCodeTemplate `yaml:"Reader"`
	ProtoCodes     *ReaderCodeTemplate `yaml:"ProtoCodes"`
	CommandMapping *CommandMapping     `yaml:"CommandMapping"`
}

type ParserResult struct {
	// 配置模板
	codeTemplates   []*CodeTemplate
	readerTemplates *ReaderCodeTemplate
	// 每个文件对应的有关键字标记的message列表
	protoMap map[string][]*ProtoMessageStructInfo // key:protoFileName
	// 所有的message
	allProto map[string][]*ProtoMessageStructInfo // key:protoFileName
}

func (this *ParserResult) GetCodeTemplate(key string) *CodeTemplate {
	for _, v := range this.codeTemplates {
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
		var structInfo *ProtoMessageStructInfo
		// struct的注释在genDecl.Doc里
		if genDecl.Doc != nil {
			keyChecker := make(map[string]struct{})
			for i, structComment := range genDecl.Doc.List {
				comment := strings.TrimPrefix(structComment.Text, "//")
				comment = strings.TrimSpace(comment)
				commentValue := ""
				normalComment := ""
				var codeTemplate *CodeTemplate
				for _, template := range parserResult.codeTemplates {
					lowerComment := strings.ToLower(comment)
					lowerKey := strings.ToLower(template.KeyComment)
					// 不区分大小写
					if lowerComment == lowerKey {
						codeTemplate = template
					} else if strings.HasPrefix(lowerComment, lowerKey) && strings.Contains(lowerComment, ":") {
						kv := strings.Split(comment, ":")
						if len(kv) == 2 && kv[1] != "" {
							codeTemplate = template
							commentValue = strings.TrimSpace(kv[1])
						}
					}
					if codeTemplate != nil && normalComment == "" {
						for j, v := range genDecl.Doc.List {
							if j != i {
								if normalComment != "" {
									normalComment += "\n"
								}
								normalComment += v.Text
							}
						}
					}
				}
				// 只处理设置了关键字的
				if codeTemplate != nil {
					// 排重
					if _, ok := keyChecker[comment]; ok {
						continue
					}
					structInfo = &ProtoMessageStructInfo{
						protoName:   path.Base(path.Clean(strings.Replace(protoCodeFile, "\\", "/", -1))),
						MessageName: typeSpec.Name.Name,
						keyComment:  codeTemplate.KeyComment,
						Value:       commentValue,
						Comment:     normalComment,
						PackageName: f.Name.Name,
						structType:  structDecl,
					}
					structInfoList := parserResult.protoMap[structInfo.protoName]
					if structInfoList == nil {
						structInfoList = make([]*ProtoMessageStructInfo, 0)
					}
					structInfoList = append(structInfoList, structInfo)
					parserResult.protoMap[structInfo.protoName] = structInfoList
					keyChecker[comment] = struct{}{}
					//println(fmt.Sprintf("%v %v key:%v value:%v", structInfo.protoName, structInfo.MessageName, structInfo.keyComment, structInfo.Value))
				}
			}
		}
		if structInfo == nil {
			normalComment := ""
			if genDecl.Doc != nil {
				for _, v := range genDecl.Doc.List {
					if normalComment != "" {
						normalComment += "\n"
					}
					normalComment += v.Text
				}
			}
			structInfo = &ProtoMessageStructInfo{
				protoName:   path.Base(path.Clean(strings.Replace(protoCodeFile, "\\", "/", -1))),
				MessageName: typeSpec.Name.Name,
				Comment:     normalComment,
				PackageName: f.Name.Name,
				structType:  structDecl,
			}
		}
		structInfoList := parserResult.allProto[structInfo.protoName]
		if structInfoList == nil {
			structInfoList = make([]*ProtoMessageStructInfo, 0)
		}
		structInfoList = append(structInfoList, structInfo)
		parserResult.allProto[structInfo.protoName] = structInfoList
	}
}

// 解析*pb.go文件
func ParseFiles(pbGoFilePattern string, codeTemplatesConfig string) {
	codeTemplates := initCodeTemplatesConfig(codeTemplatesConfig)
	parserResult := &ParserResult{
		codeTemplates:   codeTemplates.Code,
		readerTemplates: codeTemplates.Reader,
		protoMap:        map[string][]*ProtoMessageStructInfo{},
		allProto:        map[string][]*ProtoMessageStructInfo{},
	}
	files, err := filepath.Glob(pbGoFilePattern)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file count:%v", len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			continue
		}
		// It should end with ".pb.go" at a minimum.
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".pb.go") {
			continue
		}
		ParseProtoCode(file, parserResult)
	}
	// 根据proto里面的自定义tag生成代码模板
	for _, codeTemplate := range parserResult.codeTemplates {
		generateCode(parserResult, codeTemplate.KeyComment)
	}
	// 生成只读接口
	generatePbReader(parserResult)
	// 生成消息号
	generateCommandMapping(parserResult, codeTemplates.CommandMapping.OutFile)
	generateProtoCodes(parserResult, codeTemplates.ProtoCodes)
}

// 从json文件加载代码模板配置
func initCodeTemplatesConfig(config string) *CodeTemplates {
	fileData, err := os.ReadFile(config)
	if err != nil {
		panic("read config file err")
	}
	codeTemplates := &CodeTemplates{
		Code: make([]*CodeTemplate, 0),
		Reader: &ReaderCodeTemplate{
			FileFilter:    make([]string, 0),
			MessageFilter: make([]string, 0),
		},
		CommandMapping: &CommandMapping{},
	}
	err = yaml.Unmarshal(fileData, codeTemplates)
	if err != nil {
		panic(err)
	}
	autoCheckDir(&codeTemplates.Reader.OutDir)
	autoCheckDir(&codeTemplates.ProtoCodes.OutDir)
	return codeTemplates
}

func autoCheckDir(dir *string) {
	if *dir == "" {
		return
	}
	if !strings.HasSuffix(*dir, "/") && !strings.HasSuffix(*dir, "\\") {
		*dir = *dir + "/"
	}
}
