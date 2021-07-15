package scan

import "strings"

// Path 路径，最长路径是Package -> File -> Func/Struct -> Field
type Path []string

// String 打印
func (p Path) String() string {
	sb := &strings.Builder{}
	for i, s := range p {
		if i > 0 {
			sb.WriteString(" -> ")
		}
		sb.WriteString(s)
	}
	return sb.String()
}

// Clone 克隆
func (p Path) Clone() Path {
	pp := make([]string, len(p))
	copy(pp, p)
	return pp
}

// paths 生成各个组件的path，即查找路径
func (s *Scanner) paths() {
	p := Path{}
	p = append(p, s.pkg.ID)
	// 文件查找路径
	for _, file := range s.pkg.Files {
		s.filePath(p, file)
	}
}

func (s *Scanner) filePath(p Path, file *File) {
	fp := p.Clone()
	fp = append(fp, file.Name)
	file.Path = fp
	// 导入的查找路径
	for _, ip := range file.Imports {
		ifp := fp.Clone()
		ifp = append(ifp, ip.AliasName())
		ip.Path = ifp
		s.addPath(ifp, ip)
	}
	// 参数的查找路径
	for _, f := range file.Values {
		ffp := fp.Clone()
		name := f.Name
		ffp = append(ffp, name)
		f.Path = ffp
		s.addPath(ffp, f)
		// 参数的查找路径
		// for _, result := range f.Params {
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, result.Name)
		// 	result.Path = rrp
		// }
		// 结果和
		// for _, result := range f.Results {
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, result.Type)
		// 	result.Path = rrp
		// }
		// if f.Receiver != nil{
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, f.Receiver.Type)
		// 	f.Receiver.Path = rrp
		// }
	}
	// 函数的查找路径
	for _, f := range file.Funcs {
		ffp := fp.Clone()
		name := f.Name
		if f.Receiver != nil {
			name = f.Receiver.Type + "." + name
		}
		ffp = append(ffp, name)
		f.Path = ffp
		s.addPath(ffp, f)
		// 参数的查找路径
		// for _, result := range f.Params {
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, result.Name)
		// 	result.Path = rrp
		// }
		// 结果和
		// for _, result := range f.Results {
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, result.Type)
		// 	result.Path = rrp
		// }
		// if f.Receiver != nil{
		// 	rrp := ffp.Clone()
		// 	rrp = append(rrp, f.Receiver.Type)
		// 	f.Receiver.Path = rrp
		// }
	}
	for _, t := range file.Types {
		ffp := fp.Clone()
		ffp = append(ffp, t.Name)
		t.Path = ffp
		s.addPath(ffp, t)
		s.fieldPath(t, ffp)
	}
}

// fieldPath 字段查找路径
func (s *Scanner) fieldPath(t *Type, ffp Path) {
	for _, field := range t.Fields {
		fip := ffp.Clone()
		fip = append(fip, field.Name)
		field.Path = fip
		s.addPath(fip, field)
	}
}

// fieldPath 字段查找路径
func (s *Scanner) valuePath(t *Type, ffp Path) {
	for _, field := range t.Fields {
		fip := ffp.Clone()
		fip = append(fip, field.Name)
		field.Path = fip
		s.addPath(fip, field)
	}
}

func (s *Scanner) addPath(path Path, t interface{}) {
	s.pkg.PathAndTypes[path.String()] = t
}
