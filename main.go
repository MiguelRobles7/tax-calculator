package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dustin/go-humanize"
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

type DataString struct {
	Money string
    IncomeTax string
    NetPayAfterTax string
    SSS string
    PhilHealth string
    PagIbig string
    TotalContributions string
    TotalDeductions string
    NetPayAfterDeductions string
}
// Using 2023 tax table, thats what the online calculator uses but the tax table there is not updated yet
func getIncomeTax(money float64) float64 {
    if money <= 20833 {
        return 0
    } else if money <= 33332 {
        return ((money - 20833.33) * 0.1505) + 0
    } else if money <= 66666 {
        return ((money - 33333) * 0.20) + 1875
    } else if money <= 166666 {
        return ((money - 66667) * 0.25) + 8541.80
    } else if money <= 666666 {
        return ((money - 166667) * 0.30) + 33541.80 
    } else {
        return ((money - 666667) * 0.35) + 183541.80
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
    dataString := DataString{}
    http.Handle("/st/", http.StripPrefix("/st/", http.FileServer(http.Dir("st"))))
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
            data.TotalContributions = (data.SSS + data.PhilHealth + data.PagIbig + data.IncomeTax)
            data.IncomeTax = getIncomeTax(money - data.TotalContributions)
            data.NetPayAfterTax = money - data.IncomeTax
            data.TotalDeductions = data.IncomeTax + data.TotalContributions
            data.NetPayAfterDeductions =data.TotalDeductions

            dataString.Money = humanize.CommafWithDigits(money, 4)
            dataString.SSS = humanize.CommafWithDigits(data.SSS, 4)
            dataString.PhilHealth = humanize.CommafWithDigits(data.PhilHealth, 4)
            dataString.PagIbig = humanize.CommafWithDigits(data.PagIbig, 4)
            dataString.TotalContributions = humanize.CommafWithDigits(data.TotalContributions, 4)
            dataString.IncomeTax = humanize.CommafWithDigits(data.IncomeTax, 4)
            dataString.NetPayAfterTax = humanize.CommafWithDigits(data.NetPayAfterTax, 4)
            dataString.TotalDeductions = humanize.CommafWithDigits(data.TotalDeductions, 4)
            dataString.NetPayAfterDeductions = humanize.CommafWithDigits(data.NetPayAfterDeductions, 4)
        }
        tmpl.Execute(w, dataString)
    })

    http.ListenAndServe(":8080", nil)
}
