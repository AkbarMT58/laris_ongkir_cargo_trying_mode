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

	if err := trans.Origin == "" || trans.Kecamatan_destinasi == "" || trans.Berat_input == 0 || trans.Konstanta_Min_Berat == 0 || trans.Volume_input == 0 || trans.Konstanta_volume_darat_laut == 0 || trans.Konstanta_volume_udara == 0; err == true {

		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Query Param is not valid"})

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
	// var berat_final_darat int
	// var berat_final_udara int

	for _, data_all := range result.larisongkirs {
		if err := len(data_all); err == 0 {
			return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "Data Not Found"})
		}
		for i := 0; i < len(data_all); i++ {

			fmt.Println("data all i:", data_all[i])
			fmt.Println("jumlah data :", len(data_all))

			estimate_hari := data_all[i].Min_sla_hari + "-" + data_all[i].Max_sla_hari + " hari"

			tipe_pengiriman := data_all[i].Tipe

			berat_asli := berat_asli

			harga_perkg_data := data_all[i].Harga_perkg

			harga_perkg, _ := strconv.Atoi(harga_perkg_data)

			Total_Berat := Hitung_Total_Berat(tipe_pengiriman, berat_asli, trans.Konstanta_Min_Berat, int(berat_volume_metrik_daratlaut), int(berat_volume_metrik_udara))
			fmt.Println("Total_Berat:", Total_Berat)

			Total_Ongkir := Total_Berat * harga_perkg //harga per kg
			fmt.Println("Total_Ongkir:", Total_Ongkir)

			CostDesc := []CostDescription{{Cost: []Cost{{Etd: estimate_hari, Note: "", Value: float64(Total_Ongkir)}}, Description: tipe_pengiriman, Service: tipe_pengiriman}}

			join_cost_description = append(join_cost_description, CostDesc...)

		}

	}
	courier = []Courier{{Code: "laris", Costs: join_cost_description, Name: ""}}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": courier, "destination": kecamatan_dest, "origin": origin_place}})

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

func Hitung_Total_Berat(tipe_kirim string, berat_asli, k_min_berat, berat_volum_darat, berat_volum_udara int) int {

	var berat_final int
	var berat_volum int

	if tipe_kirim == "darat" {
		berat_volum = berat_volum_darat
	} else {
		berat_volum = berat_volum_udara
	}
	fmt.Println("tipe_kirim:", tipe_kirim)
	fmt.Println("berat_volum:", berat_volum)
	if berat_asli < berat_volum {

		berat_final = berat_volum

	} else {

		berat_final = berat_asli

	}

	//jika berat final < konstanta min berat

	if berat_final < k_min_berat {

		berat_final = k_min_berat

	}

	return berat_final

}

func Hitung_Total_Ongkir_Darat(tipe_kirim string, harga_perkg, berat_asli, k_min_berat, berat_volum_darat, berat_volum_udara int) int {

	var berat_final_darat int

	if berat_asli < berat_volum_darat {

		berat_final_darat = berat_volum_darat

	} else {

		berat_final_darat = berat_asli

	}

	//jika berat final < konstanta min berat

	if berat_final_darat < k_min_berat {

		berat_final_darat = k_min_berat

	}

	return berat_final_darat

}

func Hitung_Total_Ongkir_Udara(tipe_kirim string, harga_perkg, berat_asli, k_min_berat, berat_volum_darat, berat_volum_udara int) int {

	var berat_final_udara int

	if berat_asli < berat_volum_udara {

		berat_final_udara = berat_volum_udara

	} else {

		berat_final_udara = berat_asli

	}

	//jika berat final < konstanta min berat

	if berat_final_udara < k_min_berat {

		berat_final_udara = k_min_berat

	}

	return berat_final_udara

}
