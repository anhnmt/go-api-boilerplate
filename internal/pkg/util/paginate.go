package util

const (
	defaultLimitPage = 10
	defaultTotalPage = 1
)

func GetLimit(limit int) int {
	if limit <= 0 {
		limit = defaultLimitPage
	}
	return limit
}

func GetOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}

	return (page - 1) * GetLimit(limit)
}

func TotalPage(total int, pageSize int) (totalPage int) {
	if total%pageSize == 0 {
		totalPage = total / pageSize
	} else {
		totalPage = total/pageSize + 1
	}

	if totalPage == 0 {
		totalPage = defaultTotalPage
	}

	return
}
