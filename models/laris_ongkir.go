package models

import "gopkg.in/mgo.v2/bson"

type (
	Laris_ongkir struct {
		Id                  bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Origin              string        `json:"origin" bson:"origin"`
		Provinsi_destinasi  string        `json:"provinsi_destinasi" bson:"provinsi_destinasi"`
		Kota_destinasi      string        `json:"kota_destinasi" bson:"kota_destinasi"`
		Kecamatan_destinasi string        `json:"kecamatan_destinasi" bson:"kecamatan_destinasi"`
		Min_sla_hari        string        `json:"min_sla_hari" bson:"min_sla_hari"`
		Max_sla_hari        string        `json:"max_sla_hari" bson:"max_sla_hari"`
		Harga_kg            string        `json:"harga_kg" bson:"harga_kg"`
		Tipe                string        `json:"tipe" bson:"tipe"`
	}
)
