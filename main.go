package main
 
import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "os"
    "runtime"
    "strconv"
    "time"
)
 
type CurrencyDay struct {
    ID         string
    Date       time.Time
    DayNo      string
    Currencies []Currency
}
 
type Currency struct {
    Code            string
    CrossOrder      int
    Unit            int
    CurrencyNameTR  string
    CurrencyName    string
    ForexBuying     float64
    ForexSelling    float64
    BanknoteBuying  float64
    BanknoteSelling float64
    CrossRateUSD    float64
    CrossRateOther  float64
}
 
type tarih_Date struct {
    XMLName   xml.Name `xml:"Tarih_Date"`
    Tarih     string   `xml:"Tarih,attr"`
    Date      string   `xml:"Date,attr"`
    Bulten_No string   `xml:"Bulten_No,attr"`
    Currency  []currency
}
 
type currency struct {
    Kod             string `xml:"Kod,attr"`
    CrossOrder      string `xml:"CrossOrder,attr"`
    CurrencyCode    string `xml:"CurrencyCode,attr"`
    Unit            string `xml:"Unit"`
    Isim            string `xml:"Isim"`
    CurrencyName    string `xml:"CurrencyName"`
    ForexBuying     string `xml:"ForexBuying"`
    ForexSelling    string `xml:"ForexSelling"`
    BanknoteBuying  string `xml:"BanknoteBuying"`
    BanknoteSelling string `xml:"BanknoteSelling"`
    CrossRateUSD    string `xml:"CrossRateUSD"`
    CrossRateOther  string `xml:"CrossRateOther"`
}
 
func (c *CurrencyDay) GetData(CurrencyDate time.Time) {
    xDate := CurrencyDate
    t := new(tarih_Date)
    currDay := t.getData(CurrencyDate, xDate)
    for {
        if currDay == nil {
            CurrencyDate = CurrencyDate.AddDate(0, 0, -1)
            currDay := t.getData(CurrencyDate, xDate)
            if currDay != nil {
                break
            }
        } else {
            break
        }
    }
}
 
func (c *tarih_Date) getData(CurrencyDate time.Time, XDate time.Time) *CurrencyDay {
    currDay := new(CurrencyDay)
    var resp *http.Response
    var err error
    var url string
 
    currDay = new(CurrencyDay)
    url = "http://www.tcmb.gov.tr/kurlar/" + CurrencyDate.Format("200601") + "/" + CurrencyDate.Format("02012006") + ".xml"
    resp, err = http.Get(url)
 
    //fmt.Println(url)
 
    if err != nil {
        fmt.Println(err)
    } else {
 
        defer resp.Body.Close()
 
        if resp.StatusCode != http.StatusNotFound {
            tarih := new(tarih_Date)
            d := xml.NewDecoder(resp.Body)
            marshalErr := d.Decode(&tarih)
 
            if marshalErr != nil {
                log.Printf("error: %v", marshalErr)
            }
 
            c = &tarih_Date{}
            currDay.ID = XDate.Format("20060102")
            currDay.Date = XDate
            currDay.DayNo = tarih.Bulten_No
            currDay.Currencies = make([]Currency, len(tarih.Currency))
            for i, curr := range tarih.Currency {
                currDay.Currencies[i].Code = curr.CurrencyCode
                currDay.Currencies[i].CurrencyName = curr.CurrencyName
                currDay.Currencies[i].CurrencyNameTR = curr.Isim
                currDay.Currencies[i].BanknoteBuying, err = strconv.ParseFloat(curr.BanknoteBuying, 64)
                currDay.Currencies[i].BanknoteSelling, err = strconv.ParseFloat(curr.BanknoteSelling, 64)
                currDay.Currencies[i].ForexBuying, err = strconv.ParseFloat(curr.ForexBuying, 64)
                currDay.Currencies[i].ForexSelling, err = strconv.ParseFloat(curr.ForexSelling, 64)
                currDay.Currencies[i].CrossOrder, err = strconv.Atoi(curr.CrossOrder)
                currDay.Currencies[i].CrossRateOther, err = strconv.ParseFloat(curr.CrossRateOther, 64)
                currDay.Currencies[i].CrossRateUSD, err = strconv.ParseFloat(curr.CrossRateUSD, 64)
                currDay.Currencies[i].Unit, err = strconv.Atoi(curr.Unit)
            }
            //fmt.Print(currDay)     // Elde edilen veriyi konsolda göster
            //SaveJSON("currencies.json", currDay)     // Elde edilen veriyi currencies.json adında JSON formatında kaydet
 
            fmt.Println(currDay.Currencies[0].CurrencyNameTR + " / " + strconv.FormatFloat(currDay.Currencies[0].BanknoteSelling, 'E', -1, 64))
 
        } else {
            currDay = nil
        }
    }
    return currDay
}
 
func SaveJSON(fileName string, key interface{}) {
    outFile, err := os.Create(fileName)
    checkError(err)
    encoder := json.NewEncoder(outFile)
    err = encoder.Encode(key)
    checkError(err)
    outFile.Close()
}
 
func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal Error : ", err.Error())
        os.Exit(1)
    }
}
 
func main() {
 
    runtime.GOMAXPROCS(2)
    startTime := time.Now()
    CurrencyDay := new(CurrencyDay)
    CurrencyDate := time.Now()
    CurrencyDay.GetData(CurrencyDate)
 
    elepsedTime := time.Since(startTime)
    fmt.Printf("Execution Time : %s", elepsedTime)
}
