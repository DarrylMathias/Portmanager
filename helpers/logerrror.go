package helpers

func LogError(err error){
	if err != nil{
		panic(err)
	}
}