package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/AdonisVillanueva/golang-echo-mongo-api/configs"
	"github.com/AdonisVillanueva/golang-echo-mongo-api/models"
	"github.com/AdonisVillanueva/golang-echo-mongo-api/responses"
	"golang.org/x/net/context"

	"log"

	"math"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var LarisCargoCollection *mongo.Collection = configs.GetCollection(configs.DB, "laris-ongkir")
var validatecargo *validator.Validate

type array_mode struct {
	larisongkirs []larisongkir
}

type larisongkir struct {
	Origin              []string `json:"origin"`
	Provinsi_destinasi  string   `json:"provinsi_destinasi"`
	Kota_desinasi       string   `json:"kota_destinasi"`
	Kecamatan_destinasi string   `json:"kecamatan_destinasi"`
	Min_sla_hari        string   `json:"min_sla_hari"`
	Max_sla_hari        string   `json:"max_sla_hari"`
	Harga_perkg         string   `json:"harga_kg"`
	Tipe                string   `json:"tipe"`
}

type Head_Array struct {
	Data    []Data_Ongkir `json:"data"`
	Message string        `json:"message"`
	Status  string        `json:"status"`
}
type Data_Ongkir struct {
	Courier     []Courier `json:"courier"`
	Destination string    `json:"destination"`
	Origin      string    `json:"origin"`
}

type Courier struct {
	Code  string            `json:"code"`
	Costs []CostDescription `json:"costs"`
	Name  string            `json:"name"`
}

type Courier_ struct {
	Code  string            `json:"code"`
	Costs []CostDescription `json:"costs"`
	Name  string            `json:"name"`
}

type Cost struct {
	Etd   string  `json:"etd"`
	Note  string  `json:"note"`
	Value float64 `json:"value"`
}

type CostDescription struct {
	Cost        []Cost `json:"cost"`
	Description string `json:"description"`
	Service     string `json:"service"`
}

type Input_Request struct {
	Berat             int    `json:"berat"`
	Berat_Volum_Darat int    `json:"berat_volum_darat"`
	Berat_Volum_Udara int    `json:"berat_volum_udara"`
	Konstanta_berat   int    `json:"konstanta_min_berat"`
	Tipe              string `json:"tipe"`
}

// func GetAllOngkirCargo(c echo.Context) error {

func GetAllOngkir(c echo.Context) error {

	type Param_Request struct {
		Harga_by_jarak_input           string  `json:"harga_by_jarak"`
		Berat_input                    int     `json:"berat"`
		Panjang_input                  string  `json:"panjang"`
		Lebar_input                    string  `json:"lebar"`
		Tinggi_input                   string  `json:"tinggi"`
		Volume_input                   float64 `json:"volume"`
		Biaya_berat_persatuan          string  `json:"biaya_berat"`
		Biaya_volume_persatuan         string  `json:"biaya_volume"`
		Biaya_berat_limit_10_daratlaut string  `json:"biaya_berat_10_darat"`
		Biaya_berat_limit_10_udara     string  `json:"biaya_berat_10_udara"`
		Jenis_jalur_pengiriman         string  `json:"jalur_pengiriman"`
		Konstanta_volume_darat_laut    float64 `json:"konstanta_volume_darat_laut"`
		Konstanta_volume_udara         float64 `json:"konstanta_volume_udara"`
		Kecamatan_destinasi            string  `json:"kecamatan_destinasi"`
		Tipe                           string  `json:"tipe"`
		Description                    string  `json:"description"`
		Service                        string  `json:"service"`
		Notes                          string  `json:"note"`
		Name                           string  `json:"name"`
		Konstanta_Min_Berat            int     `json:"konstanta_min_berat"`
	}

	type Result_Ongkir struct {
		Total_Ongkir float64
	}

	var trans Param_Request
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(&trans)
	if err != nil {
		log.Println(err)
	}
	//deklarasi parameter
	berat_asli := trans.Berat_input
	// biaya_berat_limit_10_int_darat_laut, _ := strconv.ParseFloat(trans.Biaya_berat_limit_10_daratlaut, 64)
	// biaya_berat_limit_10_int_udara, _ := strconv.ParseFloat(trans.Biaya_berat_limit_10_udara, 64)
	konstanta_volume_darat := (trans.Konstanta_volume_darat_laut)
	konstanta_volume_udara := (trans.Konstanta_volume_udara)
	berat_volume_metrik_daratlaut := math.Ceil((trans.Volume_input) / konstanta_volume_darat)
	berat_volume_metrik_udara := math.Ceil((trans.Volume_input) / konstanta_volume_udara)
	// jenis_jalur_pengiriman := trans.Jenis_jalur_pengiriman
	// notes := trans.Notes
	// names := trans.Name

	fmt.Println("Berat Volume Metrik Darat:", math.Ceil(berat_volume_metrik_daratlaut), " Kg")
	fmt.Println("Berat Volume Metrik Udara:", math.Ceil(berat_volume_metrik_udara), " Kg")
	fmt.Println("Berat Normal:", berat_asli, " Kg")

	kecamatan_dest := string(trans.Kecamatan_destinasi)
	// biaya_berat_asli_bydata, _ := strconv.ParseFloat(v.Harga_kg, 64)
	// tipe_pengiriman := string(trans.Tipe)
	// tipe_pengiriman_darat := "darat"
	// tipe_pengiriman_udara := "udara"

	data_all_call := estimate_time_by_tipe("darat", "udara", kecamatan_dest)

	jsonData, err := json.Marshal(data_all_call)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)

	}

	fmt.Printf("json data: %s\n", jsonData)
	//fmt.Println("data asli :", data_all_call)

	dataku := []byte(`{
		"id": "363364386534336632353064343163646463393231653238",
		"origin": "Jakarta",
		"provinsi_destinasi": "Aceh",
		"kota_destinasi": "Kota Lhokseumawe",
		"kecamatan_destinasi": "Banda Sakti",
		"min_sla_hari": "5",
		"max_sla_hari": "6",
		"harga_kg": "4300",
		"tipe": "darat"
	}`)

	readdata := []array_mode{}

	result := readdata
	// result := MyData{}
	json.Unmarshal((dataku), &result)

	err = PrettyPrint(result)
	if err != nil {
		return err
	}

	//fmt.Println("Tampilkan Data All:", result)

	// harga_by_berat := 0.0
	// harga_by_volume := ((berat_volume_metrik_udara) * 4300) //m3 kubik
	// Total_Ongkir_All := (harga_by_berat + harga_by_volume)

	// Cost_darat := []CostDescription{{Cost: []Cost{{Etd: "", Note: "", Value: Total_Ongkir_All}}, Description: "", Service: ""}}
	// Cost_udara := []CostDescription{{Cost: []Cost{{Etd: "", Note: "", Value: Total_Ongkir_All}}, Description: "", Service: ""}}

	// join_group := append(Cost_darat, Cost_udara...)

	return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": data_all_call, "destination": "", "origin": ""}})

}

func hitung_berat_final(Pos Input_Request) int {

	// B_A := Pos.Berat
	// B_V_D := Pos.Berat_Volum_Darat
	// B_V_U := Pos.Berat_Volum_Udara
	// K_B := Pos.Konstanta_berat

	// if B_A > B_V_D {

	// 	B_F := B_A

	// 	fmt.Println(B_F)

	// } else {

	// 	B_F := B_V_D

	// 	fmt.Println(B_F)

	// }

	return 100

}

// func compare_konstanta_harga() int {

// }

func hitung_harga_final(c echo.Context) {

}

func estimate_time_by_tipe(td, tu, k string) interface{} {

	var dataongkir []models.Laris_ongkir

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	results, err := LarisCargoCollection.Find(ctx, bson.M{"kecamatan_destinasi": k})

	if err != nil {
		fmt.Println(err)
	}

	for results.Next(ctx) {
		var ongkirall models.Laris_ongkir
		if err = results.Decode(&ongkirall); err != nil {
			fmt.Println(err)
		}

		dataongkir = append(dataongkir, ongkirall)
	}

	var est_time_fin []string
	var data_all []interface{}

	for i, v := range dataongkir {

		if dataongkir[i].Tipe == td {

			est_time := v.Min_sla_hari + "-" + v.Max_sla_hari + " hari"

			est_time_fin = append(est_time_fin, est_time)

		}
		if dataongkir[i].Tipe == tu {

			// est_time := v.Min_sla_hari + "-" + v.Max_sla_hari + " hari"

			//est_time_fin = append(est_time_fin, est_time)

			data_all = append(data_all, dataongkir)

		}

	}

	// estimate_time := strings.Join(est_time_fin, " ")

	// fmt.Println("Hitung Waktu function:", estimate_time)

	//fmt.Println("Tampilkan Data All:", data_all)

	return data_all
}

func merge(a, b interface{}) interface{} {

	jb, err := json.Marshal(b)
	if err != nil {
		fmt.Println("Marshal error b:", err)
	}
	err = json.Unmarshal(jb, &a)
	if err != nil {
		fmt.Println("Unmarshal error b-a:", err)
	}

	return a
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}
