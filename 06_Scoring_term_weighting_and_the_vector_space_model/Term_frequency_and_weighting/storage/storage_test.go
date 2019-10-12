package storage

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	bt := InitStorage("/home/danil/Проекты/Go/Information_Retrieval/05_Index_compression/Blocked_storage/spimi/data")
	fmt.Println(Intersect(bt, "world", "do"))
}
