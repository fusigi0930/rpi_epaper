package disp

import(
	"fmt"
	"runtime"
	//"image/color"
	"os"
	"errors"
	//"reflect"

	"ggcal/log"
	"ggcal/cal"

	"gopkg.in/yaml.v3"
)

var gRoot *GObject
var gMapLayoutGlobal map[string]string

func GlobalSetting(key string) string {
	v, ok := gMapLayoutGlobal[key]
	if ok {
		return v
	}
	return ""
}

func LoadDef(file string, ext_sc *SurfaceContext) error {
	var fullname string
	gMapLayoutGlobal = make(map[string]string)

	switch runtime.GOOS {
	case "windows":
		sysdata_path := fmt.Sprintf("%s/ggcal", os.Getenv("ProgramData"))
		fullname = fmt.Sprintf("%s/%s", sysdata_path, file)
	case "linux":
		fullname = fmt.Sprintf("/etc/ggcal/%s", file)
	}

	_, err := os.Stat(fullname)
	if errors.Is(err, os.ErrNotExist) {
		log.LogService().Errorf("the layout def file %s is not exist\n", fullname)
		return err
	}

	fileData, err := os.ReadFile(fullname)
	if err != nil {
		log.LogService().Errorf("load layout def file %s failed: %v\n", fullname, err)
		return err
	}

	yamlData := string(fileData)
	var rootNode yaml.Node
	err = yaml.Unmarshal([]byte(yamlData), &rootNode)
	if err != nil {
		log.LogService().Errorf("get yaml root node failed: %v\n", err)
		return err
	}

	var screen_node *yaml.Node = &yaml.Node{}
	screen_node.Content = append(screen_node.Content, &yaml.Node{})
	for i := 0; i < len(rootNode.Content[0].Content); i += 2 {
		log.LogService().Debugf("rootname: %s\n", rootNode.Content[0].Content[i].Value)
		switch(rootNode.Content[0].Content[i].Value) {
		case "Screen":
			screen_node.Content[0].Content =
				append(screen_node.Content[0].Content, rootNode.Content[0].Content[i])
				screen_node.Content[0].Content =
				append(screen_node.Content[0].Content, rootNode.Content[0].Content[i+1])
		case "fontpath":
			log.LogService().Debugf("fontpath is: %s\n", rootNode.Content[0].Content[i+1].Value)
			gMapLayoutGlobal["fontpath"] = rootNode.Content[0].Content[i+1].Value
		}
	}

	gRoot = &GObject{SC: nil, Parent: nil}
	yamlToControl(screen_node.Content[0], gRoot)
	var sc *SurfaceContext
	if ext_sc == nil {
		log.LogService().Debugf("initial new surface...\n")
		sc = NewSurface(gRoot.GetWidth(), gRoot.GetHeight())
		log.LogService().Debugf("completed new surface...\n")
	} else {
		sc = ext_sc
	}
	defineSC(gRoot, sc)
	//gRoot.DumpParam("")

	if err := cal.LoadDef("calconfig.yaml"); err != nil {
		log.LogService().Errorf("no!!!!!!!: %v\n", err)
		return err
	}

	return nil
}

func defineSC(node GBase, sc *SurfaceContext) {
	node.SetSurfaceContext(sc)		

	for _, c := range node.GetChild() {
		defineSC(c, sc)
	}
}

func GetRootScreen() GBase {
	return gRoot
}

func yamlToControl(ynode *yaml.Node, node GBase) (error) {
	//uiType := ynode.Content[0].Value
	nodeVal := ynode.Content[1]
	childs := fillToControl(nodeVal.Content, node)
	if childs == nil {
		return nil
	}
	for _, child := range childs.Content {
		var childNode GBase
		log.LogService().Debugf("parsing %s node\n", child.Content[0].Value)
		switch (child.Content[0].Value) {
		case "Rect":
			childNode = &GRectControl{}
		case "Text":
			childNode = &GTextControl{}
		case "Image":
			childNode = &GImageControl{}
		case "CalendarGrid":
			childNode = &GCalGridControl{}
		case "LargeDay":
			childNode = &GLargeDayControl{}
		}
		childNode.Init(node)
		yamlToControl(child, childNode)
		node.AddChild(childNode)
	}
	return nil
}

func nodeToMap(nodes []*yaml.Node, leng int) map[string] string {
	m := make(map[string] string, 0)
	for i := 0;  i < leng; i += 2 {
		m[nodes[i].Value] = nodes[i+1].Value
	}
	return m
}

func fillToControl(ynode []*yaml.Node, node GBase) (*yaml.Node) {
	leng := len(ynode)
	if ynode[leng-1].Kind == yaml.SequenceNode {
		node.UpdateParam(nodeToMap(ynode, leng-2))
		return ynode[leng-1]
	} else {
		node.UpdateParam(nodeToMap(ynode, leng))
	}

	return nil
}
