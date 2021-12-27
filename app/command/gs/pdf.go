package gs

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"isbnbook/app/log"
)

type GSPdfReader struct {
	toolPath string
}

var logger = log.GetLogger()

func NewGSPdfReader(toolPath string) *GSPdfReader {
	return &GSPdfReader{toolPath}
}

func (g *GSPdfReader) GetPageCount(filePath string) (int, error) {
	//out, err := exec.Command("gs", `-q -dNOSAFER -dNODISPLAY -c "(sample.pdf) (r) file runpdfbegin pdfpagecount = quit"`).CombinedOutput()
	commandstr := fmt.Sprintf(`"(%s) (r) file runpdfbegin pdfpagecount = quit"`, filepath.Base(filePath))
	command := exec.Command(g.toolPath, "-q", "-dNOSAFER", "-dNODISPLAY", "-c", commandstr)
	// set file dir
	command.Dir = filepath.Dir(filePath)
	logger.Info(command.String())
	out, err := command.CombinedOutput()

	if err != nil {
		logger.Error("GetPageCount Error", err)
		return 0, err
	}
	num, err := strconv.Atoi(strings.TrimSpace(string(out)))
	return num, err
}

func (g *GSPdfReader) ExtractImageFiles(filePath string, outputDir string, start int, end int, prefix string) error {
	outputNames := filepath.Join(outputDir, prefix+"_image%04d.jpg")
	command := exec.Command(g.toolPath, "-q", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=jpeg", "-dJPEGQ=100",
		"-dQFactor=1.0", "-dDisplayFormat=16#30804", "-r150", "-dFirstPage="+strconv.Itoa(start), "-dLastPage="+strconv.Itoa(end), "-sOutputFile="+outputNames, filePath)
	logger.Info(command.String())
	out, err := command.CombinedOutput()

	// fmt.Println("gs " + "-q -dSAFER -dBATCH -dNOPAUSE -sDEVICE=jpeg -dJPEGQ=100 -dQFactor=1.0 -dDisplayFormat=16#30804 -r150 -dFirstPage=" + strconv.Itoa(start) + " -dLastPage=" + strconv.Itoa(end) + " -sOutputFile=" + outputNames + " " + filePath)

	if err != nil {
		logger.Error("ExtractImageFiles Error", err)
		return err
	}

	logger.Info(fmt.Sprintf("ExtractImageFiles:%s", out))

	return nil
}
