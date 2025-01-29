package service

func GetHomePageContent(c chan string) {
	s := []string{"Home", "Page"}
	func() {
		for _, v := range s {
			c <- v
		}
	}()
	close(c)
}
