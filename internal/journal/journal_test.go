package journal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"zakirullin/stuffbot/internal/userconfig"
)

func Test_AddDailyNote(t *testing.T) {
	r := require.New(t)
	now = func() time.Time {
		return time.Date(2023, 05, 30, 10, 04, 36, 0, time.UTC)
	}

	type testcase struct {
		name                string
		md                  string
		note                string
		want                string
		journalHeaderFormat string
	}

	tests := []testcase{
		{
			name: "Empty MD",
			note: "note 1",
			want: "#### 30, Tuesday\n* note 1\n",
		},
		{
			name: "New daily note",
			md:   "#### 29, Tuesday\n* note 1",
			note: "note 2",
			want: "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\n* note 2\n",
		},
		{
			name: "Append daily note",
			md:   "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\n* note 2",
			note: "note 3",
			want: "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\n* note 2\n* note 3\n",
		},

		{
			name: "Append daily note",
			md:   "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\nsome text\n* note 2",
			note: "note 3",
			want: "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\n* note 3\n\nsome text\n* note 2\n",
		},
		{
			name:                "Append daily note with custom header format",
			md:                  "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\nsome text\n* note 2",
			note:                "note 3",
			want:                "#### 29, Tuesday\n* note 1\n\n#### 30, Tuesday\n\nsome text\n* note 2\n\n#### 30.05.2023\n* note 3\n",
			journalHeaderFormat: "02.01.2006",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.journalHeaderFormat == "" {
				tt.journalHeaderFormat = userconfig.DefaultConfig.JournalHeaderFormat()
			}
			got := insertDailyNote(tt.md, tt.journalHeaderFormat, tt.note)
			r.Equal(tt.want, got)
		})
	}
}
