package df

import (
	"io"
	"strings"
	"text/template"
)

type dockerfileCfg struct {
	Hermetic bool
}

// GenerateDF generates a default dockerfile
func GenerateDF(w io.Writer, dfType string, isHermetic bool) error {
	dfc := dockerfileCfg{
		Hermetic: isHermetic,
	}

	dockerFileTmpl := getDfTmpl(strings.ToLower(dfType))

	dftmpl, err := template.New("Dockerfile").Parse(dockerFileTmpl)
	if err != nil {
		return err
	}

	err = dftmpl.Execute(w, dfc)
	if err != nil {
		return err
	}

	return nil
}

func getDfTmpl(dfType string) string {
	switch dfType {
	case "go":
		return goDfTmpl
	case "python":
		return pythonDfTmpl
	case "rust":
		return rustDfTmpl
	}
	return ""
}
