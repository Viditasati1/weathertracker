package main

import (
	"encoding/json"
	"os"
	"net/http"
	"strings"
)

type apiconfigData struct{
OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}


type weatherData struct{

	Name string `json:"name"`
   Main struct{
	Kelvin float64 `json:"temp"`       
   }`json:"main"`

}

func loadApiConfig(filename string)(apiconfigData,error){
	bytes,err:=os.ReadFile(filename) 

	if err!=nil{
		return apiconfigData{},err
	}     

	var c apiconfigData

    err=json.Unmarshal(bytes,&c);
	if err!=nil{
		return apiconfigData{},err
	}    

return  c,nil

}

func hello(w http.ResponseWriter, r* http.Request){

w.Write([]byte("hello from go\n"));
}


func query(city string)(weatherData,error){
      
apiconfig,err:=loadApiConfig(".apiconfig")

if err!=nil{
return weatherData{},err;
}
     
resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiconfig.OpenWeatherMapApiKey + "&q=" + city+ "&units=metric");

        
if err!=nil{
return weatherData{},err;
}

defer resp.Body.Close();

var d weatherData

if err:=json.NewDecoder(resp.Body).Decode(&d); err!=nil{
	return weatherData{},err
}

return  d,nil                                

}


func main(){
http.HandleFunc("/hello",hello);

http.HandleFunc("/weather/",func (w http.ResponseWriter, r* http.Request){
	city:=strings.SplitN(r.URL.Path,"/",3)[2]
	data,err:=query(city) 

	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
    w.Header().Set("Content-Type","application/json;charset=utf-8")
     json.NewEncoder(w).Encode(data);
     
})

http.ListenAndServe(":8080",nil)

}         