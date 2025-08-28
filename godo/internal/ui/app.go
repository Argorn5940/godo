package ui

import (
	"fmt"
	"godo/internal/models"
	"godo/internal/storage"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// アプリケーションの状態
type mode int

const (
	normalMode mode = iota
	inputMode
	editMode
	deleteConfirmMode
)

// アプリケーションのモデル
type Model struct {
	taskManager *models.TaskManager
	storage     *storage.TaskStorage
	cursor      int                // 選択中のタスクのインデックス
	mode        mode              // 現在のモード
	inputValue  string            // 入力中のテキスト
	editingTask int               // 編集中のタスクのインデックス
}

// 初期化関数
func NewModel() *Model {
	storage := storage.NewTaskStorage()
	tasks, err := storage.LoadTasks()
	if err != nil {
		tasks = []*models.Task{}
	}
	
	manager := models.NewTaskManager(tasks)
	
	return &Model{
		taskManager: manager,
		storage:     storage,
		cursor:      0,
		mode:        normalMode,
		inputValue:  "",
		editingTask: -1,
	}
}

// 初期化コマンド
func (m *Model) Init() tea.Cmd {
	return nil
}

// アップデート関数
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

// キー入力の処理
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case normalMode:
		return m.handleNormalMode(msg)
	case inputMode:
		return m.handleInputMode(msg)
	case editMode:
		return m.handleEditMode(msg)
	case deleteConfirmMode:
		return m.handleDeleteConfirmMode(msg)
	}
	return m, nil
}

// ノーマルモードの処理
func (m *Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	tasks := m.taskManager.GetTasks()
	
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(tasks)-1 {
			m.cursor++
		}
	case "enter":
		if len(tasks) > 0 && m.cursor < len(tasks) {
			// タスクの完了状態を切り替え
			tasks[m.cursor].Completed = !tasks[m.cursor].Completed
			m.saveToFile()
		}
	case "n":
		// 新しいタスクを追加モード
		m.mode = inputMode
		m.inputValue = ""
	case "e":
		// タスク編集モード
		if len(tasks) > 0 && m.cursor < len(tasks) {
			m.mode = editMode
			m.editingTask = m.cursor
			m.inputValue = tasks[m.cursor].Title
		}
	case "d":
		// タスク削除確認モード
		if len(tasks) > 0 && m.cursor < len(tasks) {
			m.mode = deleteConfirmMode
		}
	}
	return m, nil
}

// 入力モードの処理
func (m *Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.inputValue) != "" {
			m.taskManager.AddTask(strings.TrimSpace(m.inputValue))
			m.saveToFile()
		}
		m.mode = normalMode
		m.inputValue = ""
	case "esc":
		m.mode = normalMode
		m.inputValue = ""
	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
	default:
		if len(msg.String()) == 1 && len(m.inputValue) < 30 {
			m.inputValue += msg.String()
		}
	}
	return m, nil
}

// 編集モードの処理
func (m *Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.inputValue) != "" {
			tasks := m.taskManager.GetTasks()
			if m.editingTask < len(tasks) {
				tasks[m.editingTask].Title = strings.TrimSpace(m.inputValue)
				m.saveToFile()
			}
		}
		m.mode = normalMode
		m.inputValue = ""
		m.editingTask = -1
	case "esc":
		m.mode = normalMode
		m.inputValue = ""
		m.editingTask = -1
	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
	default:
		if len(msg.String()) == 1 && len(m.inputValue) < 30 {
			m.inputValue += msg.String()
		}
	}
	return m, nil
}

// 削除確認モードの処理
func (m *Model) handleDeleteConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		// タスクを削除
		tasks := m.taskManager.GetTasks()
		if m.cursor < len(tasks) {
			// スライスから要素を削除
			copy(tasks[m.cursor:], tasks[m.cursor+1:])
			tasks = tasks[:len(tasks)-1]
			m.taskManager = models.NewTaskManager(tasks)
			
			// カーソル位置を調整
			if m.cursor >= len(tasks) && len(tasks) > 0 {
				m.cursor = len(tasks) - 1
			}
			if len(tasks) == 0 {
				m.cursor = 0
			}
			
			m.saveToFile()
		}
		m.mode = normalMode
	case "n", "esc":
		m.mode = normalMode
	}
	return m, nil
}

// ファイルに保存
func (m *Model) saveToFile() {
	m.storage.SaveTasks(m.taskManager.GetTasks())
}

// ビュー関数
func (m *Model) View() string {
	var s strings.Builder
	
	// スタイル定義
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)
		
	completedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")) // 緑
		
	incompleteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")) // オレンジ
		
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Padding(0, 1)
		
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)
		
	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	// ヘッダー
	completed, total := m.taskManager.GetStats()
	header := fmt.Sprintf("📄 Godo - タスク管理    完了: %d | 未完了: %d", completed, total-completed)
	s.WriteString(headerStyle.Render(header))
	s.WriteString("\n\n")

	// タスクリスト
	tasks := m.taskManager.GetTasks()
	if len(tasks) == 0 {
		s.WriteString("タスクがありません。'n'で新しいタスクを追加してください。\n")
	} else {
		for i, task := range tasks {
			var status string
			var taskStyle lipgloss.Style
			
			if task.Completed {
				status = "✓"
				taskStyle = completedStyle
			} else {
				status = "○"
				taskStyle = incompleteStyle
			}
			
			taskLine := fmt.Sprintf("%s %s", status, task.Title)
			dateLine := dateStyle.Render(fmt.Sprintf("    作成: %s | 更新: %s", 
				task.CreatedAt.Format("2006-01-02 15:04"),
				task.UpdatedAt.Format("2006-01-02 15:04")))
			
			if i == m.cursor {
				s.WriteString(selectedStyle.Render(taskStyle.Render(taskLine)))
				s.WriteString("\n")
				s.WriteString(selectedStyle.Render(dateLine))
			} else {
				s.WriteString(taskStyle.Render(taskLine))
				s.WriteString("\n")
				s.WriteString(dateLine)
			}
			s.WriteString("\n\n")
		}
	}

	// モード別の表示
	switch m.mode {
	case inputMode:
		s.WriteString("\n新しいタスクを入力してください (30文字まで):\n")
		s.WriteString(fmt.Sprintf("> %s", m.inputValue))
		s.WriteString("\n\nEnter: 追加 | Esc: キャンセル")
		
	case editMode:
		s.WriteString("\nタスクを編集してください (30文字まで):\n")
		s.WriteString(fmt.Sprintf("> %s", m.inputValue))
		s.WriteString("\n\nEnter: 保存 | Esc: キャンセル")
		
	case deleteConfirmMode:
		if len(tasks) > 0 && m.cursor < len(tasks) {
			s.WriteString(fmt.Sprintf("\n'%s' を削除しますか？\n", tasks[m.cursor].Title))
			s.WriteString("y: はい | n: いいえ")
		}
		
	default:
		// フッター（操作説明）
		footer := "操作: Enter=完了切替 | n=追加 | e=編集 | d=削除 | ↑↓=選択 | q=終了"
		s.WriteString("\n")
		s.WriteString(footerStyle.Render(footer))
	}

	return s.String()
}

// TUIアプリケーションを開始する関数
func RunApp() error {
	model := NewModel()
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}