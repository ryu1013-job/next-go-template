package migrate

import (
	"bufio"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

func Run(ctx context.Context, db *sql.DB, dir string) error {
	files, _ := filepath.Glob(filepath.Join(dir, "*.sql"))
	for _, f := range files {
		b, _ := os.ReadFile(f)
		up := extractUp(string(b))
		if strings.TrimSpace(up) == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, up); err != nil {
			return err
		}
	}
	return nil
}
func extractUp(s string) string {
	sc := bufio.NewScanner(strings.NewReader(s))
	var b strings.Builder
	in := false
	for sc.Scan() {
		l := sc.Text()
		if strings.Contains(l, "+migrate Up") {
			in = true
			continue
		}
		if strings.Contains(l, "+migrate Down") {
			break
		}
		if in {
			b.WriteString(l + "\n")
		}
	}
	if b.Len() == 0 {
		return s
	}
	return b.String()
}
