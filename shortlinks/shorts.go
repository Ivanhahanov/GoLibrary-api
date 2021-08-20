package shortlinks

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Shorter struct {
	ShortName    string `json:"short_name"`
	OriginalName string `json:"original_name"`
	VisitCount   int    `json:"visit_count"`
}

var ctx = context.Background()
var linkDB = redis.Client{}
var countDB = redis.Client{}


func InitRedisConnection() {
	linkDB = *redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "sOmE_sEcUrE_pAsS", // no password set
		DB:       0,  // use default DB
	})
	countDB = *redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "sOmE_sEcUrE_pAsS", // no password set
		DB:       1,  // use default DB
	})

}

func CreateShortLink(originalLink string) (string, error) {
	shortLink := GenerateShorLink(4)
	err := linkDB.Set(ctx, shortLink, originalLink, 0).Err()
	if err != nil {
		return "", err
	}
	err = countDB.Set(ctx, shortLink, 0, 0).Err()
	if err != nil {
		return "", err
	}
	return shortLink, nil
}

func GetOriginalLink(shortLink string) (originalLink string, err error) {
	originalLink, err = linkDB.Get(ctx, shortLink).Result()
	if err != nil {
		return "", err
	}
	return originalLink, nil
}

func GetAllDocuments() ([]*Shorter, error) {

	var shorters []*Shorter
	keys, err := linkDB.Keys(ctx, "*").Result()
	if err != nil{
		return nil, err
	}
	for _, key := range keys{
		orig, _ := linkDB.Get(ctx, key).Result()
		count, _ := countDB.Get(ctx, key).Int()
		shorters = append(shorters, &Shorter{
			ShortName: key,
			OriginalName: orig,
			VisitCount: count,
		})
	}
	return shorters, nil
}

func WriteVisit(shortLink string) error {
	countDB.Incr(ctx, shortLink)
	return nil
}
