package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	Version string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		Version: "V0.0.1",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// FileInfo represents information about a PDF file
type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Pages int    `json:"pages"`
}

// ProcessResult represents the result of processing a PDF file
type ProcessResult struct {
	File    string `json:"file"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SelectFiles allows the user to select multiple PDF files
func (a *App) SelectFiles() ([]FileInfo, error) {
	files := make([]FileInfo, 0)

	// 获取用户选择的文件
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择PDF文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF文件 (*.pdf)",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("选择文件失败: %v", err)
	}

	// 处理每个选择的文件
	for _, path := range filePaths {
		// 确保路径是绝对路径
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("获取绝对路径失败: %v, 原路径: %s\n", err, path)
			continue
		}

		// 获取文件信息
		pageCount, err := api.PageCountFile(absPath)
		if err != nil {
			// 记录详细错误信息
			fmt.Printf("获取PDF页数失败: %v, 文件: %s\n", err, absPath)

			// 尝试使用默认页数
			fmt.Printf("使用默认页数1继续处理文件\n")
			pageCount = 1
		}

		files = append(files, FileInfo{
			Name:  filepath.Base(absPath),
			Path:  absPath,
			Pages: pageCount,
		})
	}

	return files, nil
}

// parsePageRange 解析页码范围字符串并返回页码列表和是否为连续范围
func parsePageRange(pageRange string) ([]int, bool, error) {
	pages := []int{}
	isContinuous := true
	ranges := strings.Split(pageRange, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			bounds := strings.Split(r, "-")
			if len(bounds) != 2 {
				return nil, false, fmt.Errorf("无效的页码范围: %s", r)
			}
			start, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, false, fmt.Errorf("无效的起始页码: %s", bounds[0])
			}
			end, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, false, fmt.Errorf("无效的结束页码: %s", bounds[1])
			}
			for i := start; i <= end; i++ {
				pages = append(pages, i)
			}
		} else {
			page, err := strconv.Atoi(r)
			if err != nil {
				return nil, false, fmt.Errorf("无效的页码: %s", r)
			}
			pages = append(pages, page)
			isContinuous = false
		}
	}
	return pages, isContinuous, nil
}

// ExtractPages extracts specified pages from PDF files
func (a *App) ExtractPages(filePaths []string, pageRange string) []ProcessResult {
	results := make([]ProcessResult, 0)

	// 解析页码范围
	pages, isContinuous, err := parsePageRange(pageRange)
	if err != nil {
		result := ProcessResult{
			File:    "",
			Success: false,
			Message: fmt.Sprintf("解析页码范围失败: %v", err),
		}
		results = append(results, result)
		return results
	}

	// 将 pages 转换为 []string
	pageStrs := make([]string, len(pages))
	for i, p := range pages {
		pageStrs[i] = strconv.Itoa(p)
	}

	for _, path := range filePaths {
		result := ProcessResult{
			File:    filepath.Base(path),
			Success: true,
		}

		// 清理和准备文件路径
		dir := filepath.Dir(path)
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(filepath.Base(path), ext)

		// 创建output文件夹
		outDir := filepath.Join(dir, "output")
		err := os.MkdirAll(outDir, 0755)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("创建输出目录失败: %v", err)
			results = append(results, result)
			continue
		}

		// 清理文件名，移除特殊字符
		safeName := strings.Map(func(r rune) rune {
			switch {
			case r >= 'a' && r <= 'z',
				r >= 'A' && r <= 'Z',
				r >= '0' && r <= '9',
				r == '-', r == '_':
				return r
			case r >= '\u4e00' && r <= '\u9fa5': // 保留中文字符
				return r
			default:
				return '_'
			}
		}, base)

		if isContinuous {
			// 连续页码提取
			outFileName := fmt.Sprintf("%s_p%d-%d%s", safeName, pages[0], pages[len(pages)-1], ext)
			outPath := filepath.Join(outDir, outFileName)
			// 确保输入路径为绝对路径
			absPath, err := filepath.Abs(path)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("获取绝对路径失败: %v", err)
				results = append(results, result)
				continue
			}
			absOutPath, err := filepath.Abs(outPath)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("获取输出文件绝对路径失败: %v", err)
				results = append(results, result)
				continue
			}
			// 提取并合并页码
			err = api.ExtractPagesFile(absPath, outDir, pageStrs, nil)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("提取页面失败: %v", err)
				results = append(results, result)
				continue
			}
			err = api.MergeCreateFile([]string{filepath.Join(outDir, fmt.Sprintf("%s_page_%d%s", safeName, pages[0], ext))}, absOutPath, false, nil)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("合并提取页面失败: %v", err)
			} else {
				result.Message = fmt.Sprintf("已保存到: %s", outFileName)
			}
		} else {
			// 非连续页码提取
			for _, pageNumber := range pages {
				outFileName := fmt.Sprintf("%s_page_%d%s", safeName, pageNumber, ext)
				// 确保输入路径为绝对路径
				absPath, err := filepath.Abs(path)
				if err != nil {
					result.Success = false
					result.Message = fmt.Sprintf("获取绝对路径失败: %v", err)
					results = append(results, result)
					continue
				}
				// 修改：直接使用输出目录outDir作为第二个参数，而非absOutPath
				err = api.ExtractPagesFile(absPath, outDir, []string{strconv.Itoa(pageNumber)}, nil)
				if err != nil {
					result.Success = false
					result.Message = fmt.Sprintf("提取页面 %d 失败: %v", pageNumber, err)
				} else {
					result.Message = fmt.Sprintf("已保存到: %s", outFileName)
				}
				results = append(results, result)
			}
		}
	}

	return results
}
