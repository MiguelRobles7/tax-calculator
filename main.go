package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Data struct {
	IsGet bool
	Money float64
    IncomeTax float64
    NetPayAfterTax float64
    SSS float64
    PhilHealth float64
    PagIbig float64
    TotalContributions float64
    TotalDeductions float64
    NetPayAfterDeductions float64
}

// di pa toh accurate @miguelle
func getIncomeTax(money float64) float64 {
    if money <= 20833 {
        return 0
    } else if money <= 33332 {
        return (money * 0.2)
    } else if money <= 66666 {
        return (money - 2500) * 0.25
    } else if money <= 166666 {
        return (money - 10833.33) * 0.3
    } else if money <= 666666 {
        return (money - 40833.33) * 0.32
    } else {
        return (money - 200833.33) * 0.35
    }
}
// @izabelle
func getSSS(money float64) float64{
    var minLimit float64 = 4250
    var maxLimit float64 = 29750
    var sss float64 = 180
    if money >= maxLimit{
        sss=1350.00
    } else if money >= minLimit {
	    for i:= minLimit; i<=money; i+=500 {
	    	sss+=22.5
	    }
    } 
    return sss
}

// @jmse - done (will update after clarify w/ sir)
// PhilHeatlh calculator uses 2023 Employed Membership: 2.25% rate 
func getPhilHealth(money float64) float64{
    var philHealth float64 = 0
    if money <= 10000.00 {
        philHealth = 225.00
    } else if money >= 10000.01 && money <= 89999.99 {
        philHealth = money * 0.0225
    } else if money >= 90000.00 {
        philHealth = 4050.00
    }
    return philHealth
}

// @miguelle - done - will ask sir if capped at 100 talaga cause thats what the online tool does
func getPagIbig(money float64) float64{
    var pagIbig float64 = 0
    if money <= 1500{
        pagIbig = money * 0.01
    } else{
        pagIbig = money * 0.02
    }
    if pagIbig > 100 {
        pagIbig = 100
    }
    return pagIbig
}

func main() {
    tmpl := template.Must(template.ParseFiles("index.html"))
    data := Data{IsGet: false, Money: 0}

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            println("GET")
            data.IsGet = true
            data.Money = float64(0)
        } else {
            println("POST")
            var money, _ = strconv.ParseFloat(r.FormValue("money"), 64)
            data.IsGet = false
            data.Money = money
            data.SSS = getSSS(money)
            data.PhilHealth = getPhilHealth(money)
            data.PagIbig = getPagIbig(money)
            data.TotalContributions = data.SSS + data.PhilHealth + data.PagIbig
            data.IncomeTax = getIncomeTax(money)
            data.NetPayAfterTax = money - data.IncomeTax
            data.TotalDeductions = data.IncomeTax + data.TotalContributions
            data.NetPayAfterDeductions = data.NetPayAfterTax - data.TotalDeductions
        }
        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
