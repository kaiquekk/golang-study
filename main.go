package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "goApp"
	password = "pass@1234"
	dbname   = "cotasHist"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/generateChart/bar/{year}", handleGenerateBarChart).Methods("GET")
	router.HandleFunc("/addData", handleAddDataFromFile).Methods("POST")
	fmt.Println("Servidor esperando requisições na porta 8089.")
	log.Fatal(http.ListenAndServe(":8089", router))
}

func handleAddDataFromFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.WriteHeader(http.StatusOK)
	readFile(r)
}

func handleGenerateBarChart(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnectionString)
	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)
	fmt.Println("Conexão com a base de dados estabelecida com sucesso!")

	generateBarChart(db, convertToInt(params["year"]), w)
}

func readFile(r *http.Request) {
	file, handler, err := r.FormFile("datafile")
	checkError(err)
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	checkError(err)
	defer f.Close()

	io.Copy(f, file)

	psqlConnectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnectionString)
	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)
	fmt.Println("Conexão com a base de dados estabelecida com sucesso!")

	uploadedFile, err := os.Open(handler.Filename)
	checkError(err)

	fileScanner := bufio.NewScanner(uploadedFile)

	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	firstLine := fileScanner.Text()[2:31]
	row := db.QueryRow("SELECT arquivo FROM goapp.arquivos_processados WHERE arquivo = $1", firstLine)

	var arquivo string
	err = row.Scan(&arquivo)
	if err == sql.ErrNoRows {
		_, err = db.Exec(`INSERT INTO goapp.arquivos_processados (arquivo) VALUES ($1)`, firstLine)
		checkError(err)

		for fileScanner.Scan() {
			line := fileScanner.Text()
			if line[0:2] != "99" {
				insertQuery := `INSERT INTO goapp.dados_cota (data,
				precoabertura, precomax, precomin, precoultimo,
				precomedio, totalnegocios, qtdtitulos, voltitulos,
				nomeempresa, codnegocio)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
				_, err = db.Exec(insertQuery,
					formatDate(line[2:10]),
					convertToFloat(line[56:67], line[67:69]),
					convertToFloat(line[69:80], line[80:82]),
					convertToFloat(line[82:93], line[93:95]),
					convertToFloat(line[108:119], line[119:121]),
					convertToFloat(line[95:106], line[106:108]),
					convertToInt(line[147:152]),
					convertToInt64(line[152:170]),
					convertToInt64(line[170:188]),
					line[27:39],
					line[12:24])

				checkError(err)
			}
		}
		file.Close()
	}
	fmt.Println("Processo de leitura de arquivo finalizado.")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func formatDate(dateText string) string {
	return fmt.Sprintf("%s-%s-%s", dateText[:4], dateText[4:6], dateText[6:8])
}

func convertToInt(text string) int {
	var value, err = strconv.Atoi(text)
	checkError(err)
	return value
}

func convertToInt64(text string) int64 {
	var value, err = strconv.ParseInt(text, 10, 64)
	checkError(err)
	return value
}

func convertToFloat(wholeNumber string, decimalNumber string) float64 {
	var value, err = strconv.ParseFloat(fmt.Sprintf("%s.%s", wholeNumber, decimalNumber), 64)
	checkError(err)
	return value
}

func generateBarChart(db *sql.DB, year int, w http.ResponseWriter) {
	rows, err := db.Query(`SELECT TO_CHAR(ss.numeroMes, 'Mon') AS mes, ss.precoMax FROM (
	SELECT DATE_TRUNC('month', dados_cota.data) AS numeroMes, MAX(dados_cota.precomax) as precoMax 
	FROM goapp.dados_cota WHERE EXTRACT(YEAR FROM dados_cota.data) = $1 GROUP BY numeroMes ORDER BY numeroMes) AS ss;`, year)
	checkError(err)
	defer rows.Close()
	var months []string
	prices := make([]opts.BarData, 0)

	for rows.Next() {
		var mes string
		var precomax int64
		err = rows.Scan(&mes, &precomax)
		checkError(err)
		months = append(months, mes)
		prices = append(prices, opts.BarData{Value: precomax})
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Preço Máximo dentre Todas as Ações por Mês",
			Subtitle: "Maior preço dentre todas as ações que foram negociadas em cada mês",
		}),
		charts.WithLegendOpts(opts.Legend{Right: "30%"}),
	)
	bar.SetXAxis(months).
		AddSeries("Preço Máximo", prices)
	// f, _ := os.Create("chart.html")
	bar.Render(w)
}
