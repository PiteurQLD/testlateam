package scrapper

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type aggregatedNickname struct {
	Nickname string
	UIDs     []string
}

func NicknamesToUIDs(pCtx context.Context, nicks []string) (map[string][]string, error) {
	idC := make(chan aggregatedNickname)
	g, ctx := errgroup.WithContext(pCtx)
	for _, n := range nicks {
		n := n
		g.Go(func() error {
			res, err := SearchNickname(ctx, n)
			if err == nil {
				nIds := make([]string, 0, len(res.Data))
				for _, data := range res.Data {
					nIds = append(nIds, data.EncryptedUID)
				}
				idC <- aggregatedNickname{Nickname: n, UIDs: nIds}
			}
			return err
		})
	}
	go func() {
		g.Wait()
		close(idC)
	}()
	aRes := make(map[string][]string, len(nicks))
	for id := range idC {
		aRes[id.Nickname] = id.UIDs
	}
	return aRes, g.Wait()
}
