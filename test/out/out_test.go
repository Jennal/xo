package out

import (
	"database/sql"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testxo?charset=utf8&parseTime=True&loc=Local")
	assert.NotNil(t, db)
	assert.Nil(t, err)
	defer db.Close()

	u := &User{
		Property: &UserProperty{
			Nickname: sql.NullString{
				String: "nickname",
				Valid:  true,
			},
		},
		Prop2: &Prop2{
			Field1: "abc",
			Field2: 123,
		},
		Name: sql.NullString{
			String: "name",
			Valid:  true,
		},
		Age: sql.NullInt64{
			Int64: 23,
			Valid: true,
		},

		Properties: []*UserProperty{
			&UserProperty{
				Nickname: sql.NullString{
					String: "name1",
					Valid:  true,
				},
			},
			&UserProperty{
				Nickname: sql.NullString{
					String: "name2",
					Valid:  true,
				},
			},
			&UserProperty{
				Nickname: sql.NullString{
					String: "name3",
					Valid:  true,
				},
			},
		},
	}

	err = u.Save(db)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testxo?charset=utf8&parseTime=True&loc=Local")
	assert.NotNil(t, db)
	assert.Nil(t, err)
	defer db.Close()

	u, err := UserByID(db, 8)
	assert.NotNil(t, u)
	assert.Nil(t, err)

	spew.Dump(u)
}

func TestUpdate(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testxo?charset=utf8&parseTime=True&loc=Local")
	assert.NotNil(t, db)
	assert.Nil(t, err)
	defer db.Close()

	u, err := UserByID(db, 10)
	assert.NotNil(t, u)
	assert.Nil(t, err)

	u.Prop2.Field1 = "xyz"
	u.Property.Nickname.String = "jennal"

	u.Save(db)

	spew.Dump(u)
}

func TestDelete(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testxo?charset=utf8&parseTime=True&loc=Local")
	assert.NotNil(t, db)
	assert.Nil(t, err)
	defer db.Close()

	u, err := UserByID(db, 10)
	assert.NotNil(t, u)
	assert.Nil(t, err)

	u.Delete(db)

	spew.Dump(u)
}
