package bot

import (
	"strconv"
	"strings"
	"testing"
)

func TestGenerateProductListMarkerAddsPaginationButtons(t *testing.T) {
	markup := GenerateProductListMarker(ProductListPage{
		Products: []ShopProduct{
			{ID: "p1", Name: "Product 1", Price: 10},
		},
		CategoryID: "c1",
		Page:       2,
		PageSize:   ProductListPageSize,
		Total:      11,
		TotalPages: 3,
	})

	if len(markup.InlineKeyboard) != 2 {
		t.Fatalf("expected product row and pagination row, got %d rows", len(markup.InlineKeyboard))
	}

	paginationRow := markup.InlineKeyboard[1]
	if len(paginationRow) != 3 {
		t.Fatalf("expected previous, current, and next buttons, got %d buttons", len(paginationRow))
	}

	expectedPages := []int{1, 2, 3}
	for i, expectedPage := range expectedPages {
		callbackData := paginationRow[i].CallbackData
		if callbackData == nil {
			t.Fatalf("expected callback data on pagination button %d", i)
		}
		if !strings.Contains(*callbackData, "page="+strconv.Itoa(expectedPage)) {
			t.Fatalf("expected callback %q to contain page=%d", *callbackData, expectedPage)
		}
	}
}
