package util

func GetOffset(page, limit int) int {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return (page - 1) * limit
}

func TotalPage(total int, pageSize int) (totalPage int) {
	if total%pageSize == 0 {
		totalPage = total / pageSize
	} else {
		totalPage = total/pageSize + 1
	}

	if totalPage == 0 {
		totalPage = 1
	}

	return
}
