package service


func (S *Service) GetHistory(from, to int) map[string]map[string]string {
	messages,err := S.Database.HistoryMessages(from,to)
	if err != nil {
		// err message
	}
	return messages
}