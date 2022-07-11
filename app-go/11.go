package main


func myFunction() (string){
	return "hello"
}

func b() (string, error){
	for i := 0; i < 8; i++ {
       	myFunction()
        }
        return "Success",nil
}
func main() {
	b()
}

 


