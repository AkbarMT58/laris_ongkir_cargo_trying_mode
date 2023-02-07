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

type Test struct {
	Code  string `json:"code"`
	Costs string `json:"costs"`
	Name  string `json:"name"`
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
	notes := trans.Notes
	names := trans.Name

	fmt.Println("Berat Volume Metrik Darat:", math.Ceil(berat_volume_metrik_daratlaut), " Kg")
	fmt.Println("Berat Volume Metrik Udara:", math.Ceil(berat_volume_metrik_udara), " Kg")
	fmt.Println("Berat Normal:", berat_asli, " Kg")

	kecamatan_dest := string(trans.Kecamatan_destinasi)
	// tipe_pengiriman := string(trans.Tipe)
	tipe_pengiriman_darat := "darat"
	tipe_pengiriman_udara := "udara"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var dataongkir []models.Laris_ongkir
	defer cancel()

	results, err := LarisCargoCollection.Find(ctx, bson.M{"kecamatan_destinasi": kecamatan_dest})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var ongkirall models.Laris_ongkir
		if err = results.Decode(&ongkirall); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		dataongkir = append(dataongkir, ongkirall)
	}

	for i, v := range dataongkir {

		var time_darat []string
		var time_udara []string

		if dataongkir[i].Tipe == "darat" {

			Estimate_Time_Darat := v.Min_sla_hari + "-" + v.Max_sla_hari + " hari"

			time_darat = append(time_darat, Estimate_Time_Darat)

		}

		if dataongkir[i].Tipe == "udara" {

			Estimate_Time_Udara := v.Min_sla_hari + "-" + v.Max_sla_hari + " hari"

			time_udara = append(time_udara, Estimate_Time_Udara)

		}

		// estimate_darat_laut := strings.Join(time_darat, " ")
		// estimate_darat_udara := strings.Join(time_udara, " ")

		fmt.Println("Estimate Time darat:", time_darat)
		fmt.Println("Estimate Time udara:", time_udara)
		fmt.Println("Harga Per Kg:", v.Harga_kg)

		//fmt.Println("Hitung Berat Aktual:", hitung_berat_final(Input_Request{berat_asli, int(berat_volume_metrik_daratlaut), int(berat_volume_metrik_udara), trans.Konstanta_Min_Berat, tipe_pengiriman}))

		// deklarasi by data
		biaya_berat_asli_bydata, _ := strconv.ParseFloat(v.Harga_kg, 64)

		Estimate_Time := v.Min_sla_hari + "-" + v.Max_sla_hari + " hari"

		if berat_asli <= trans.Konstanta_Min_Berat {

			if berat_asli <= trans.Konstanta_Min_Berat {

				//kondisi if berat volum > berat asli
				if int(berat_volume_metrik_daratlaut) > berat_asli {

					if int(berat_volume_metrik_daratlaut) <= trans.Konstanta_Min_Berat {

						if tipe_pengiriman_darat == "darat" {

							harga_by_volume := 0.0
							harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

							//Total perhitungan semua aspek

							Total_Ongkir_All := (harga_by_berat + harga_by_volume)

							group := []Courier{{

								Code:  "laris",
								Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
								Name:  names,
							}}

							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

							fmt.Println(" 1 a total ongkir darat(kondisi metrik berat < konstanta berat 10):", Total_Ongkir_All)
						}

					} else {

						if tipe_pengiriman_darat == v.Tipe {

							//darat dan laut

							harga_by_berat := 0.0
							harga_by_volume := ((berat_volume_metrik_daratlaut) * biaya_berat_asli_bydata) //m3 kubik
							Total_Ongkir_All := (harga_by_berat + harga_by_volume)

							group := []Courier{{
								Code:  "laris",
								Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
								Name:  names,
							},
							}

							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

							fmt.Println("1 b total ongkir volum metrik saat > 10 kg di darat:", Total_Ongkir_All)

						}

						if tipe_pengiriman_udara == v.Tipe {

							//darat dan laut

							harga_by_berat := 0.0
							harga_by_volume := ((berat_volume_metrik_udara) * biaya_berat_asli_bydata) //m3 kubik
							Total_Ongkir_All := (harga_by_berat + harga_by_volume)

							Cost_darat := []CostDescription{{Cost: []Cost{{Etd: "", Note: "", Value: Total_Ongkir_All}}, Description: "", Service: ""}}
							Cost_udara := []CostDescription{{Cost: []Cost{{Etd: "", Note: "", Value: Total_Ongkir_All}}, Description: "", Service: ""}}

							join_group := append(Cost_darat, Cost_udara...)

							Courier := []Courier{{

								Code:  "laris",
								Costs: []CostDescription{{}},
								Name:  names,
							}}

							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": join_group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

							fmt.Println("1 b total ongkir volum metrik saat > 10 kg di udara:", Courier)

						}

					}

				}

				if int(berat_volume_metrik_udara) > berat_asli {

					if int(berat_volume_metrik_udara) <= trans.Konstanta_Min_Berat {

						if tipe_pengiriman_udara == v.Tipe {

							harga_by_volume := 0.0
							harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

							//Total perhitungan semua aspek

							Total_Ongkir_All := (harga_by_berat + harga_by_volume)

							group := []Courier{{

								Code:  "laris",
								Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
								Name:  names,
							}}

							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

							fmt.Println(" 1 c total ongkir volum metrik udara(kondisi berat volum metrik udara < konstanta berat 10):", Total_Ongkir_All)

						}

					}

					// else {

					// 	if tipe_pengiriman_udara == v.Tipe {

					// 		//udara

					// 		harga_by_berat := 0.0
					// 		harga_by_volume := ((berat_volume_metrik_udara) * biaya_berat_asli_bydata) //m3 kubik
					// 		Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					// 		group := []Courier{{

					// 			Code:  "laris",
					// 			Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
					// 			Name:  names,
					// 		}}

					// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					// 		fmt.Println("1 d total ongkir saat > 10 kg di udara:", Total_Ongkir_All)

					// 	}
					// }

				}

				//jika berat asli > volume metrik

				if int(berat_volume_metrik_daratlaut) < berat_asli {

					if tipe_pengiriman_darat == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 2 a total ongkir darat laut(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

					}

					if int(berat_volume_metrik_udara) < berat_asli {

						if tipe_pengiriman_udara == v.Tipe {

							harga_by_volume := 0.0
							harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

							//Total perhitungan semua aspek

							Total_Ongkir_All := (harga_by_berat + harga_by_volume)

							group := []Courier{{

								Code:  "laris",
								Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
								Name:  names,
							}}

							c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

							fmt.Println(" 2 b total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

						}

					}

				}

				if int(berat_volume_metrik_daratlaut) == berat_asli {

					if tipe_pengiriman_darat == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 2 c total ongkir darat laut(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

					}

				}

				if int(berat_volume_metrik_udara) == berat_asli {

					if tipe_pengiriman_udara == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 2 d total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

					}

				}

			} else {
				//batas

				if tipe_pengiriman_darat == v.Tipe {

					harga_by_volume := 0.0
					harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

					//Total perhitungan semua aspek

					Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					group := []Courier{{

						Code:  "laris",
						Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
						Name:  names,
					}}

					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					fmt.Println(" 2 e total ongkir darat laut(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

				}

				if tipe_pengiriman_udara == v.Tipe {

					harga_by_volume := 0.0
					harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

					//Total perhitungan semua aspek

					Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					group := []Courier{{

						Code:  "laris",
						Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
						Name:  names,
					}}

					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					fmt.Println(" 2 f total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

				}

			}

			//batas kondisi

		} else {

			//kondisi if berat volum > berat asli
			if int(berat_volume_metrik_daratlaut) > berat_asli {

				if int(berat_volume_metrik_daratlaut) <= trans.Konstanta_Min_Berat {

					if tipe_pengiriman_darat == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 3 a total ongkir darat(kondisi metrik berat < konstanta berat 10):", Total_Ongkir_All)
					}

				} else {

					if tipe_pengiriman_darat == v.Tipe {

						//darat dan laut

						harga_by_berat := 0.0
						harga_by_volume := ((berat_volume_metrik_daratlaut) * biaya_berat_asli_bydata) //m3 kubik
						Total_Ongkir_All := (harga_by_berat + harga_by_volume)
						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println("3 b total ongkir volum metrik saat > 10 kg di darat:", Total_Ongkir_All)

					}

				}

			}

			if int(berat_volume_metrik_udara) > berat_asli {

				if int(berat_volume_metrik_udara) <= trans.Konstanta_Min_Berat {

					if tipe_pengiriman_udara == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 3 c total ongkir volum metrik udara(kondisi berat volum metrik udara < konstanta berat 10):", Total_Ongkir_All)

					}

				} else {

					if tipe_pengiriman_udara == v.Tipe {

						//udara

						harga_by_berat := 0.0
						harga_by_volume := ((berat_volume_metrik_udara) * biaya_berat_asli_bydata) //m3 kubik
						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println("3 d total ongkir saat > 10 kg di udara:", Total_Ongkir_All)

					}
				}

			}

			if int(berat_volume_metrik_daratlaut) <= berat_asli {

				if tipe_pengiriman_darat == v.Tipe {

					harga_by_volume := 0.0
					harga_by_berat := (float64(berat_asli) * biaya_berat_asli_bydata) //harga per kg

					//Total perhitungan semua aspek

					Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					group := []Courier{{

						Code:  "laris",
						Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
						Name:  names,
					}}

					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					fmt.Println(" 4 a total ongkir darat laut(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

				}

				if int(berat_volume_metrik_udara) <= berat_asli {

					if tipe_pengiriman_udara == v.Tipe {

						harga_by_volume := 0.0
						harga_by_berat := (float64(berat_asli) * biaya_berat_asli_bydata) //harga per kg

						//Total perhitungan semua aspek

						Total_Ongkir_All := (harga_by_berat + harga_by_volume)

						group := []Courier{{

							Code:  "laris",
							Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
							Name:  names,
						}}

						c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

						fmt.Println(" 4 b total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

					}

				}

			}

			if int(berat_volume_metrik_daratlaut) == berat_asli {

				if tipe_pengiriman_darat == v.Tipe {

					harga_by_volume := 0.0
					harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

					//Total perhitungan semua aspek

					Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					group := []Courier{{

						Code:  "laris",
						Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
						Name:  names,
					}}

					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					fmt.Println(" 4 c total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

				}

			}

			if int(berat_volume_metrik_udara) == berat_asli {

				if tipe_pengiriman_udara == v.Tipe {

					harga_by_volume := 0.0
					harga_by_berat := (float64(trans.Konstanta_Min_Berat) * biaya_berat_asli_bydata) //harga per kg

					//Total perhitungan semua aspek

					Total_Ongkir_All := (harga_by_berat + harga_by_volume)

					group := []Courier{{

						Code:  "laris",
						Costs: []CostDescription{{Cost: []Cost{{Etd: Estimate_Time, Note: notes, Value: Total_Ongkir_All}}, Description: v.Tipe, Service: v.Tipe}},
						Name:  names,
					}}

					c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"courier": group, "destination": v.Kecamatan_destinasi, "origin": v.Origin}})

					fmt.Println(" 4 d total ongkir udara(kondisi berat < konstanta berat 10):", Total_Ongkir_All)

				}

			}

		}

	}

	//batas berat volume metrik

	return nil
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
