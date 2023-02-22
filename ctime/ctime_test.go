package ctime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	now = func() time.Time {
		// 2022-11-08T14:52:33Z
		return time.Date(2022, time.November, 8, 14, 52, 33, 0, time.UTC)
	}
	defer func() { now = time.Now }()
	m.Run()
}

func TestNew(t *testing.T) {
	ct := Now()
	assert.Equal(t, 2022, ct.Year())
	assert.Equal(t, time.November, ct.Month())
	assert.Equal(t, 8, ct.Day())
	assert.Equal(t, 14, ct.Hour())
	assert.Equal(t, 52, ct.Minute())
	assert.Equal(t, 33, ct.Second())
	assert.Equal(t, 0, ct.Nanosecond())
	assert.Equal(t, time.UTC, ct.Location())

	std := now()
	ct = New(std)
	assert.Equal(t, 2022, ct.Year())
	assert.Equal(t, time.November, ct.Month())
	assert.Equal(t, 8, ct.Day())
	assert.Equal(t, 14, ct.Hour())
	assert.Equal(t, 52, ct.Minute())
	assert.Equal(t, 33, ct.Second())
	assert.Equal(t, 0, ct.Nanosecond())
	assert.Equal(t, time.UTC, ct.Location())

}

func TestFormat(t *testing.T) {
	ct := Now()
	wanted := "2022-11-08T14:52:33Z"

	assert.Equal(t, wanted, ct.String())

	writer := bytes.NewBuffer(nil)
	fmt.Fprintf(writer, "%v", ct)
	assert.Equal(t, wanted, writer.String())
}

func TestIn(t *testing.T) {
	ct := Now()

	assert.Equal(t, ct, ct.In(nil))

	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)

	ct = ct.In(loc)

	assert.Equal(t, 2022, ct.Year())
	assert.Equal(t, time.November, ct.Month())
	assert.Equal(t, 8, ct.Day())
	assert.Equal(t, 14+8, ct.Hour())
	assert.Equal(t, 52, ct.Minute())
	assert.Equal(t, 33, ct.Second())
	assert.Equal(t, 0, ct.Nanosecond())
	assert.Equal(t, loc, ct.Location())

	wanted := "2022-11-08T22:52:33+08:00"

	assert.Equal(t, wanted, ct.String())

	writer := bytes.NewBuffer(nil)
	fmt.Fprintf(writer, "%v", ct)
	assert.Equal(t, wanted, writer.String())
}

func TestJson(t *testing.T) {
	type testStruct struct {
		Name      string `json:"name"`
		CreatedAt CTime  `json:"created_at"`
	}

	s := &testStruct{
		Name:      "abc",
		CreatedAt: Now(),
	}

	b, err := json.Marshal(s)
	require.NoError(t, err)

	wanted := `{"name":"abc","created_at":"2022-11-08T14:52:33Z"}`
	assert.Equal(t, wanted, string(b))

	input := `{"name":"abc","created_at":"2022-11-08T22:52:33+08:00"}`
	result := new(testStruct)
	err = json.Unmarshal([]byte(input), result)
	require.NoError(t, err)
	assert.Equal(t, "abc", result.Name)
	assert.Equal(t, 2022, result.CreatedAt.Year())
	assert.Equal(t, time.November, result.CreatedAt.Month())
	assert.Equal(t, 8, result.CreatedAt.Day())
	assert.Equal(t, 14, result.CreatedAt.Hour())
	assert.Equal(t, 52, result.CreatedAt.Minute())
	assert.Equal(t, 33, result.CreatedAt.Second())
	assert.Equal(t, 0, result.CreatedAt.Nanosecond())
	assert.Equal(t, time.UTC, result.CreatedAt.Location())

	empty := &testStruct{
		Name:      "xxx",
		CreatedAt: CTime{},
	}

	b, err = json.Marshal(empty)
	require.NoError(t, err)

	wanted = `{"name":"xxx","created_at":null}`
	assert.Equal(t, wanted, string(b))

	input = `{"name":"test","created_at":null}`
	result = new(testStruct)
	err = json.Unmarshal([]byte(input), result)
	require.NoError(t, err)
	assert.Equal(t, true, result.CreatedAt.IsZero())

	errInput := `{"name":"test","created_at":"2022-11-0822:52:33"}`
	result = new(testStruct)
	err = json.Unmarshal([]byte(errInput), result)
	require.Error(t, err)
}

type testDBStruct struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt CTime  `gorm:"column:created_at"`
}

func (testDBStruct) TableName() string {
	return "t_test"
}

func TestSql(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(t, err)

	// write
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `t_test` (`name`,`created_at`) VALUES (?,?)")).
		WithArgs("abc", Now()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	testWrite := &testDBStruct{
		Name:      "abc",
		CreatedAt: Now(),
	}
	err = db.Debug().Create(testWrite).Error
	require.NoError(t, err)

	// read
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_test` WHERE id = ? ORDER BY `t_test`.`id` LIMIT 1")).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "created_at"}).
				AddRow(1, "abc", time.Date(2022, time.November, 8, 14, 52, 33, 0, time.UTC)),
		)

	testRead := new(testDBStruct)
	err = db.Debug().First(testRead, "id = ?", 1).Error
	require.NoError(t, err)

	assert.Equal(t, uint64(1), testRead.ID)
	assert.Equal(t, "abc", testRead.Name)
	assert.Equal(t, Now(), testRead.CreatedAt)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
