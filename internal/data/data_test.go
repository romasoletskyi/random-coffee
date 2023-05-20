package data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	d, err := CreateRawFormDatabase(ctx)
	db := CreateFormDatabase(d)
	require.NoError(t, err)

	form := UserForm{"name", "email", "", "bio", "", mapInfo{0.0, 0.0, 2000}, make([][]int, 0), make([]string, 0)}
	db.AddUserForm(ctx, form)
	require.NoError(t, err)
}
