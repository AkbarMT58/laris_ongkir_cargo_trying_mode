package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	larisongkirs [][]larisongkir
}

type larisongkir struct {
	Origin              string `json:"origin"`
	Provinsi_destinasi  string `json:"provinsi_destinasi"`
	Kota_destinasi      string `json:"kota_destinasi"`
	Kecamatan_destinasi string `json:"kecamatan_destinasi"`
	Min_sla_hari        string `json:"min_sla_hari"`
	Max_sla_hari        string `json:"max_sla_hari"`
	Harga_perkg         string `json:"harga_kg"`
	Tipe                string `json:"tipe"`
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
		Origin                         string  `json:"origin"`
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
	origin_place := trans.Origin
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

	data_all_call := estimate_time_by_tipe("darat", "udara", kecamatan_dest, origin_place)

	jsonData, err := json.Marshal(data_all_call)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)

	}

	//fmt.Println("data asli :", data_all_call)

	readdata := array_mode{}
	result := readdata

	if err := json.Unmarshal((jsonData), &result.larisongkirs); err != nil {
		fmt.Println(err)
		panic(1)
	}

	// var join_ongkir []interface{}

	var join_cost_description []CostDescription
	var courier []Courier
	var berat_final_darat int
	var berat_final_udara int

	for _, data_all := range result.larisongkirs {

		for i := 0; i < len(data_all); i++ {

			fmt.Println("data all i:", data_all[i])
			fmt.Println("jumlah data :", len(data_all))

			estimate_darat_laut := data_all[0].Min_sla_hari + "-" + data_all[0].Max_sla_hari + " hari"
			estimate_darat_udara := data_all[1].Min_sla_hari + "-" + data_all[1].Max_sla_hari + " hari"

			// destinasi := data_all[0].Kecamatan_destinasi
			// origin := data_all[0].Origin

			tipe_pengiriman_darat := data_all[0].Tipe
			tipe_pengiriman_udara := data_all[1].Tipe
			berat_asli := berat_asli

			harga_perkg_data_darat := data_all[0].Harga_perkg
			harga_perkg_darat, _ := strconv.Atoi(harga_perkg_data_darat)

			harga_perkg_data_udara := data_all[1].Harga_perkg
			harga_perkg_udara, _ := strconv.Atoi(harga_perkg_data_udara)

			//Total_Ongkir := Hitung_Total_Ongkir(tipe_pengiriman, harga_perkg_int, berat_asli, int(konstanta_volume_darat), int(konstanta_volume_udara), trans.Konstanta_Min_Berat, int(berat_volume_metrik_daratlaut), int(berat_volume_metrik_udara))
			// Total_Ongkir := 20000000

			//kondisi if

			//jika berat volum darat > berat asli

			if berat_asli < int(berat_volume_metrik_daratlaut) {

				berat_final_darat = int(berat_volume_metrik_daratlaut)

			} else {

				berat_final_darat = berat_asli

			}

			//jika berat volum udara > berat asli

			if berat_asli < int(berat_volume_metrik_udara) {

				berat_final_udara = int(berat_volume_metrik_udara)

			} else {

				berat_final_udara = berat_asli

			}

			//jika berat final < konstanta min berat

			if berat_final_darat < trans.Konstanta_Min_Berat {

				berat_final_darat = trans.Konstanta_Min_Berat

			}

			//jika berat volume udara < konstanta min berat
			if berat_final_udara < trans.Konstanta_Min_Berat {

				berat_final_udara = trans.Konstanta_Min_Berat

			}

			//darat

			harga_by_berat_darat := berat_final_darat * harga_perkg_darat //harga per kg
			Total_Ongkir_All_Darat := harga_by_berat_darat

			//udara

			harga_by_berat_udara := berat_final_udara * harga_perkg_udara //harga per kg
			Total_Ongkir_All_Udara := harga_by_berat_udara

			// Cost_darat := []Courier{{Code: "laris", Costs: []CostDescription{{Cost: []Cost{{Etd: estimate_darat_laut, Note: "", Value: float64(Total_Ongkir_All_Darat)}}, Description: tipe_pengiriman_darat, Service: tipe_pengiriman_darat}}, Name: ""}}
			// Cost_udara := []Courier{{Code: "laris", Costs: []CostDescription{{Cost: []Cost{{Etd: estimate_darat_udara, Note: "", Value: float64(Total_Ongkir_All_Udara)}}, Description: tipe_pengiriman_udara, Service: tipe_pengiriman_udara}}, Name: ""}}

			Cost_darat := []CostDescription{{Cost: []Cost{{Etd: estimate_darat_laut, Note: "", Value: float64(Total_Ongkir_All_Darat)}}, Description: tipe_pengiriman_darat, Service: tipe_pengiriman_darat}}
			Cost_udara := []CostDescription{{Cost: []Cost{{Etd: estimate_darat_udara, Note: "", Value: float64(Total_Ongkir_All_Udara)}}, Description: tipe_pengiriman_udara, Service: tipe_pengiriman_udara}}

			join_cost_description = append(Cost_darat, Cost_udara...)

			//batas kondisi

		}

	}
	courier = []Courier{{Code: "laris", Costs: join_cost_description, Name: ""}}
	return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": courier, "destination": kecamatan_dest, "origin": origin_place}})

}

func estimate_time_by_tipe(td, tu, k, o string) interface{} {

	var dataongkir []models.Laris_ongkir

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	results, err := LarisCargoCollection.Find(ctx, bson.M{"kecamatan_destinasi": k, "origin": o})

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

	var data_all []interface{}

	data_all = append(data_all, dataongkir)

	return data_all
}

// func Hitung_Total_Ongkir(tipe_kirim string, harga_perkg, berat_asli, k_d, k_u, k_b, berat_volum_darat, berat_volum_udara int) int {

// 	var total_ongkir int

// 	if berat_asli < k_b {

// 		tot_ongkir := (berat_asli) * harga_perkg

// 		total_ongkir = append(total_ongkir, tot_ongkir)

// 	}

// 	return total_ongkir

// }

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
