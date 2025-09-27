package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// キーメッセージ作成ヘルパー
func key(s string) tea.KeyMsg {
	// 1文字はRunes、特殊キーはTypeのみで十分
	if len(s) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

// モデルにキー入力を順に送る
func sendKeys(m *Model, keys ...string) *Model {
	for _, k := range keys {
		mm, _ := m.Update(key(k))
		m = mm.(*Model)
	}
	return m
}

func TestNewModel_InitState(t *testing.T) {
	// 実ファイルに触れないよう環境を隔離
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	if m == nil {
		t.Fatalf("NewModel returned nil")
	}
	if m.mode != normalMode {
		t.Fatalf("expected normalMode, got %v", m.mode)
	}
	if m.cursor != 0 || m.inputValue != "" || m.editingTask != -1 {
		t.Fatalf("unexpected initial fields: cursor=%d input=%q editing=%d", m.cursor, m.inputValue, m.editingTask)
	}
	if len(m.taskManager.GetTasks()) != 0 {
		t.Fatalf("expected no tasks initially")
	}
}

func TestAddTask_viaInputMode(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	// 'n' -> 入力モード、"abc"入力 -> Enter
	m = sendKeys(m, "n", "a", "b", "c", "enter")

	tasks := m.taskManager.GetTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "abc" || tasks[0].Completed {
		t.Fatalf("unexpected task: %+v", tasks[0])
	}
	if m.mode != normalMode || m.inputValue != "" {
		t.Fatalf("should return to normal mode with empty input")
	}
}

func TestToggleTaskWithEnter(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	m = sendKeys(m, "n", "x", "enter")
	if m.taskManager.GetTasks()[0].Completed {
		t.Fatalf("should start incomplete")
	}
	// normal モードで Enter でトグル
	m = sendKeys(m, "enter")
	if !m.taskManager.GetTasks()[0].Completed {
		t.Fatalf("expected completed after toggle")
	}
}

func TestEditTaskTitle(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	m = sendKeys(m, "n", "x", "enter")
	// 編集へ -> 既存タイトルに追記/置換を簡単にするためBackspaceで消して新規文字列
	m = sendKeys(m, "e")
	// 既存タイトル長だけBackspace
	for range m.inputValue {
		m = sendKeys(m, "backspace")
	}
	m = sendKeys(m, "n", "e", "w", "enter")

	if got := m.taskManager.GetTasks()[0].Title; got != "new" {
		t.Fatalf("expected title 'new', got %q", got)
	}
	if m.mode != normalMode || m.inputValue != "" || m.editingTask != -1 {
		t.Fatalf("should exit edit mode and reset fields")
	}
}

func TestDeleteTaskWithConfirmation(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	m = sendKeys(m, "n", "a", "enter", "n", "b", "enter")
	if len(m.taskManager.GetTasks()) != 2 {
		t.Fatalf("setup failed: need 2 tasks")
	}
	// 先頭にカーソルのまま d -> y で削除
	m = sendKeys(m, "d")
	if m.mode != deleteConfirmMode {
		t.Fatalf("expected deleteConfirmMode")
	}
	m = sendKeys(m, "y")
	if len(m.taskManager.GetTasks()) != 1 {
		t.Fatalf("expected 1 task after delete, got %d", len(m.taskManager.GetTasks()))
	}
	if m.mode != normalMode {
		t.Fatalf("should return to normal mode after delete")
	}
}

func TestCursorNavigationBounds(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	m = sendKeys(m, "n", "a", "enter", "n", "b", "enter")
	// 下へ
	m = sendKeys(m, "down")
	if m.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", m.cursor)
	}
	// 末尾でさらにdownしても変化しない
	m = sendKeys(m, "down")
	if m.cursor != 1 {
		t.Fatalf("cursor should stay at last index")
	}
	// upで戻る、先頭ではこれ以上上がらない
	m = sendKeys(m, "up", "up")
	if m.cursor != 0 {
		t.Fatalf("cursor should stay at 0, got %d", m.cursor)
	}
}

func TestViewShowsStates(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	m := NewModel()
	out := m.View()
	if !strings.Contains(out, "タスクがありません") {
		t.Fatalf("empty state should be shown, got: %s", out)
	}
	// タスク追加後の表示
	m = sendKeys(m, "n", "a", "enter")
	out = m.View()
	if !strings.Contains(out, "a") {
		t.Fatalf("task title should be in view: %s", out)
	}
	if !strings.Contains(out, "操作: Enter=完了切替") {
		t.Fatalf("footer help should be present")
	}
}


