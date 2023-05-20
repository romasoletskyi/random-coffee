package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/romasoletskyi/random-coffee/internal/app"
	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/romasoletskyi/random-coffee/internal/user"
	"github.com/sirupsen/logrus"
)

func main() {
	file, db := app.Initialize("connect-log", data.CreateRawPairDatabase)
	defer file.Close()
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pairdb := data.CreatePairDatabase(db)
	forms, err := pairdb.GetPairs(ctx)
	if err != nil {
		logrus.Error(err)
	}

	for _, form := range forms {
		fmt.Printf("%v %v <---> %v %v\n", form.Left.Name, form.Left.Email, form.Right.Name, form.Right.Email)
	}
	fmt.Println("Do you want to proceed (y/n)?")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	if answer[0] == 'y' {
		var wg sync.WaitGroup

		for i, form := range forms {
			wg.Add(2)
			go func(i int, f data.PairForm) {
				defer wg.Done()
				time.Sleep(time.Duration(2*i) * time.Second)
				err := user.SendInvitationMail(f)
				if err != nil {
					logrus.Error(err)
				}
			}(i, form)
			go func(i int, f data.PairForm) {
				defer wg.Done()
				time.Sleep(time.Duration(2*i+1) * time.Second)
				err := user.SendInvitationMail(data.ReversePairForm(f))
				if err != nil {
					logrus.Error(err)
				}
			}(i, form)
		}

		wg.Wait()
	} else {
		fmt.Println("Aborting")
	}
}
